package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// CommsOperation provides operations to manage the cloudCommunications singleton.
type CommsOperation struct {
    Entity
    // Unique Client Context string. Max limit is 256 chars.
    clientContext *string
    // The result information. Read-only.
    resultInfo ResultInfoable
    // The status property
    status *OperationStatus
}
// NewCommsOperation instantiates a new commsOperation and sets the default values.
func NewCommsOperation()(*CommsOperation) {
    m := &CommsOperation{
        Entity: *NewEntity(),
    }
    odataTypeValue := "#microsoft.graph.commsOperation";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateCommsOperationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateCommsOperationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.addLargeGalleryViewOperation":
                        return NewAddLargeGalleryViewOperation(), nil
                    case "#microsoft.graph.cancelMediaProcessingOperation":
                        return NewCancelMediaProcessingOperation(), nil
                    case "#microsoft.graph.inviteParticipantsOperation":
                        return NewInviteParticipantsOperation(), nil
                    case "#microsoft.graph.muteParticipantOperation":
                        return NewMuteParticipantOperation(), nil
                    case "#microsoft.graph.playPromptOperation":
                        return NewPlayPromptOperation(), nil
                    case "#microsoft.graph.recordOperation":
                        return NewRecordOperation(), nil
                    case "#microsoft.graph.startHoldMusicOperation":
                        return NewStartHoldMusicOperation(), nil
                    case "#microsoft.graph.stopHoldMusicOperation":
                        return NewStopHoldMusicOperation(), nil
                    case "#microsoft.graph.subscribeToToneOperation":
                        return NewSubscribeToToneOperation(), nil
                    case "#microsoft.graph.unmuteParticipantOperation":
                        return NewUnmuteParticipantOperation(), nil
                    case "#microsoft.graph.updateRecordingStatusOperation":
                        return NewUpdateRecordingStatusOperation(), nil
                }
            }
        }
    }
    return NewCommsOperation(), nil
}
// GetClientContext gets the clientContext property value. Unique Client Context string. Max limit is 256 chars.
func (m *CommsOperation) GetClientContext()(*string) {
    return m.clientContext
}
// GetFieldDeserializers the deserialization information for the current model
func (m *CommsOperation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["clientContext"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetClientContext)
    res["resultInfo"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateResultInfoFromDiscriminatorValue , m.SetResultInfo)
    res["status"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseOperationStatus , m.SetStatus)
    return res
}
// GetResultInfo gets the resultInfo property value. The result information. Read-only.
func (m *CommsOperation) GetResultInfo()(ResultInfoable) {
    return m.resultInfo
}
// GetStatus gets the status property value. The status property
func (m *CommsOperation) GetStatus()(*OperationStatus) {
    return m.status
}
// Serialize serializes information the current object
func (m *CommsOperation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("clientContext", m.GetClientContext())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("resultInfo", m.GetResultInfo())
        if err != nil {
            return err
        }
    }
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err = writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetClientContext sets the clientContext property value. Unique Client Context string. Max limit is 256 chars.
func (m *CommsOperation) SetClientContext(value *string)() {
    m.clientContext = value
}
// SetResultInfo sets the resultInfo property value. The result information. Read-only.
func (m *CommsOperation) SetResultInfo(value ResultInfoable)() {
    m.resultInfo = value
}
// SetStatus sets the status property value. The status property
func (m *CommsOperation) SetStatus(value *OperationStatus)() {
    m.status = value
}
