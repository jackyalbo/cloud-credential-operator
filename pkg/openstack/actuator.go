/*
Copyright 2019 The OpenShift Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package openstack

import (
	"context"
	"fmt"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"

	configv1 "github.com/openshift/api/config/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	log "github.com/sirupsen/logrus"

	minterv1 "github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential/v1"
	"github.com/openshift/cloud-credential-operator/pkg/operator/constants"
	actuatoriface "github.com/openshift/cloud-credential-operator/pkg/operator/credentialsrequest/actuator"
	"github.com/openshift/cloud-credential-operator/pkg/operator/utils"
)

type OpenStackActuator struct {
	Client client.Client
	Codec  *minterv1.ProviderCodec
}

// NewOpenStackActuator creates a new OpenStack actuator.
func NewOpenStackActuator(client client.Client) (*OpenStackActuator, error) {
	codec, err := minterv1.NewCodec()
	if err != nil {
		log.WithError(err).Error("error creating OpenStack codec")
		return nil, fmt.Errorf("error creating OpenStack codec: %v", err)
	}

	return &OpenStackActuator{
		Codec:  codec,
		Client: client,
	}, nil
}

// Exists checks if the credentials currently exist.
func (a *OpenStackActuator) Exists(ctx context.Context, cr *minterv1.CredentialsRequest) (bool, error) {
	logger := a.getLogger(cr)
	logger.Debug("running Exists")
	var err error

	existingSecret := &corev1.Secret{}
	err = a.Client.Get(context.TODO(), types.NamespacedName{Namespace: cr.Spec.SecretRef.Namespace, Name: cr.Spec.SecretRef.Name}, existingSecret)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Debug("target secret does not exist")
			return false, nil
		}
		return false, err
	}

	logger.Debug("target secret exists")
	return true, nil
}

// Create the credentials.
func (a *OpenStackActuator) Create(ctx context.Context, cr *minterv1.CredentialsRequest) error {
	logger := a.getLogger(cr)
	logger.Debug("running create")
	return a.sync(ctx, cr)
}

// Update the credentials to the provided definition.
func (a *OpenStackActuator) Update(ctx context.Context, cr *minterv1.CredentialsRequest) error {
	logger := a.getLogger(cr)
	logger.Debug("running update")
	return a.sync(ctx, cr)
}

func (a *OpenStackActuator) sync(ctx context.Context, cr *minterv1.CredentialsRequest) error {
	logger := a.getLogger(cr)
	logger.Debug("running sync")

	credentialsRootSecret, err := a.GetCredentialsRootSecret(ctx, cr)
	if err != nil {
		logger.WithError(err).Error("issue with cloud credentials secret")
		return err
	}

	clouds, err := GetRootCloudCredentialsSecretData(credentialsRootSecret, logger)
	if err != nil {
		logger.WithError(err).Error("issue with cloud credentials secret")
		return &actuatoriface.ActuatorError{
			ErrReason: minterv1.InsufficientCloudCredentials,
			Message:   err.Error(),
		}
	}

	logger.Debugf("provisioning secret")
	existingSecret, err := a.loadExistingSecret(cr)
	if err != nil {
		return err
	}

	err = a.syncCredentialSecret(cr, clouds, existingSecret, "", logger)
	if err != nil {
		msg := "error creating/updating secret"
		logger.WithError(err).Error(msg)
		return &actuatoriface.ActuatorError{
			ErrReason: minterv1.CredentialsProvisionFailure,
			Message:   fmt.Sprintf("%v: %v", msg, err),
		}
	}

	return nil
}

func (a *OpenStackActuator) syncCredentialSecret(cr *minterv1.CredentialsRequest, clouds string, existingSecret *corev1.Secret, userPolicy string, logger log.FieldLogger) error {
	sLog := logger.WithFields(log.Fields{
		"targetSecret": fmt.Sprintf("%s/%s", cr.Spec.SecretRef.Namespace, cr.Spec.SecretRef.Name),
		"cr":           fmt.Sprintf("%s/%s", cr.Namespace, cr.Name),
	})

	if existingSecret == nil || existingSecret.Name == "" {
		if clouds == "" {
			msg := "new access key secret needed but no key data provided"
			sLog.Error(msg)
			return &actuatoriface.ActuatorError{
				ErrReason: minterv1.CredentialsProvisionFailure,
				Message:   msg,
			}
		}
		sLog.Info("creating secret")
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cr.Spec.SecretRef.Name,
				Namespace: cr.Spec.SecretRef.Namespace,
				Annotations: map[string]string{
					minterv1.AnnotationCredentialsRequest: fmt.Sprintf("%s/%s", cr.Namespace, cr.Name),
				},
			},
			Data: map[string][]byte{RootOpenStackCredsSecretKey: []byte(clouds)},
		}

		err := a.Client.Create(context.TODO(), secret)
		if err != nil {
			sLog.WithError(err).Error("error creating secret")
			return err
		}
		sLog.Info("secret created successfully")
		return nil
	}

	// Update the existing secret:
	sLog.Debug("updating secret")
	origSecret := existingSecret.DeepCopy()
	if existingSecret.Annotations == nil {
		existingSecret.Annotations = map[string]string{}
	}
	existingSecret.Annotations[minterv1.AnnotationCredentialsRequest] = fmt.Sprintf("%s/%s", cr.Namespace, cr.Name)
	if clouds != "" {
		existingSecret.Data[RootOpenStackCredsSecretKey] = []byte(clouds)
	}

	if !reflect.DeepEqual(existingSecret, origSecret) {
		sLog.Info("target secret has changed, updating")
		err := a.Client.Update(context.TODO(), existingSecret)
		if err != nil {
			msg := "error updating secret"
			sLog.WithError(err).Error(msg)
			return &actuatoriface.ActuatorError{
				ErrReason: minterv1.CredentialsProvisionFailure,
				Message:   msg,
			}
		}
	} else {
		sLog.Debug("target secret unchanged")
	}

	return nil
}

func (a *OpenStackActuator) loadExistingSecret(cr *minterv1.CredentialsRequest) (*corev1.Secret, error) {
	logger := a.getLogger(cr)

	// Check if the credentials secret exists, if not we need to inform the syncer to generate a new one:
	existingSecret := &corev1.Secret{}
	err := a.Client.Get(context.TODO(), types.NamespacedName{Namespace: cr.Spec.SecretRef.Namespace, Name: cr.Spec.SecretRef.Name}, existingSecret)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Debug("secret does not exist")
		} else {
			return nil, err
		}
	}

	if _, ok := existingSecret.Data[RootOpenStackCredsSecretKey]; !ok {
		logger.Warning("secret did not have expected key: " + RootOpenStackCredsSecretKey)
	} else {
		logger.Debug("found clouds.yaml in existing secret")
	}

	return existingSecret, nil
}

// GetCredentialsRootSecretLocation returns the namespace and name where the parent credentials secret is stored.
func (a *OpenStackActuator) GetCredentialsRootSecretLocation() types.NamespacedName {
	return types.NamespacedName{Namespace: constants.CloudCredSecretNamespace, Name: constants.OpenStackCloudCredsSecretName}
}

func (a *OpenStackActuator) GetCredentialsRootSecret(ctx context.Context, cr *minterv1.CredentialsRequest) (*corev1.Secret, error) {
	logger := a.getLogger(cr)
	cloudCredSecret := &corev1.Secret{}
	if err := a.Client.Get(ctx, a.GetCredentialsRootSecretLocation(), cloudCredSecret); err != nil {
		msg := "unable to fetch root cloud cred secret"
		logger.WithError(err).Error(msg)
		return nil, &actuatoriface.ActuatorError{
			ErrReason: minterv1.CredentialsProvisionFailure,
			Message:   fmt.Sprintf("%v: %v", msg, err),
		}
	}
	return cloudCredSecret, nil
}

// Delete credentials
func (a *OpenStackActuator) Delete(ctx context.Context, cr *minterv1.CredentialsRequest) error {
	logger := a.getLogger(cr)
	logger.Debug("running delete")

	existingSecret := &corev1.Secret{}
	err := a.Client.Get(context.TODO(), types.NamespacedName{Namespace: cr.Spec.SecretRef.Namespace, Name: cr.Spec.SecretRef.Name}, existingSecret)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Debug("target secret does not exist")
			return nil
		}
		return err
	}

	err = a.Client.Delete(context.TODO(), existingSecret)
	if err != nil {
		return err
	}

	return nil
}

func (a *OpenStackActuator) getLogger(cr *minterv1.CredentialsRequest) log.FieldLogger {
	return log.WithFields(log.Fields{
		"actuator": "openstack",
		"cr":       fmt.Sprintf("%s/%s", cr.Namespace, cr.Name),
	})
}

// Upgradeable returns a ClusterOperator status condition for the upgradeable type
// if the system is considered not upgradeable. Otherwise, return nil as the default
// value is for things to be upgradeable.
func (a *OpenStackActuator) Upgradeable(mode operatorv1.CloudCredentialsMode) *configv1.ClusterOperatorStatusCondition {
	return utils.UpgradeableCheck(a.Client, mode, a.GetCredentialsRootSecretLocation())
}
