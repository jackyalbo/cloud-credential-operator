package models
import (
    "errors"
)
// Provides operations to manage the collection of agreementAcceptance entities.
type DetectedAppPlatformType int

const (
    // Default. Set to unknown when platform cannot be determined.
    UNKNOWN_DETECTEDAPPPLATFORMTYPE DetectedAppPlatformType = iota
    // Indicates that the platform of the detected application is Windows.
    WINDOWS_DETECTEDAPPPLATFORMTYPE
    // Indicates that the platform of the detected application is Windows Mobile.
    WINDOWSMOBILE_DETECTEDAPPPLATFORMTYPE
    // Indicates that the platform of the detected application is Windows Holographic.
    WINDOWSHOLOGRAPHIC_DETECTEDAPPPLATFORMTYPE
    // Indicates that the platform of the detected application is iOS.
    IOS_DETECTEDAPPPLATFORMTYPE
    // Indicates that the platform of the detected application is macOS.
    MACOS_DETECTEDAPPPLATFORMTYPE
    // Indicates that the platform of the detected application is ChromeOS.
    CHROMEOS_DETECTEDAPPPLATFORMTYPE
    // Indicates that the platform of the detected application is Android open source project.
    ANDROIDOSP_DETECTEDAPPPLATFORMTYPE
    // Indicates that the platform of the detected application is Android device administrator.
    ANDROIDDEVICEADMINISTRATOR_DETECTEDAPPPLATFORMTYPE
    // Indicates that the platform of the detected application is Android work profile.
    ANDROIDWORKPROFILE_DETECTEDAPPPLATFORMTYPE
    // Indicates that the platform of the detected application is Android dedicated and fully managed.
    ANDROIDDEDICATEDANDFULLYMANAGED_DETECTEDAPPPLATFORMTYPE
)

func (i DetectedAppPlatformType) String() string {
    return []string{"unknown", "windows", "windowsMobile", "windowsHolographic", "ios", "macOS", "chromeOS", "androidOSP", "androidDeviceAdministrator", "androidWorkProfile", "androidDedicatedAndFullyManaged"}[i]
}
func ParseDetectedAppPlatformType(v string) (interface{}, error) {
    result := UNKNOWN_DETECTEDAPPPLATFORMTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_DETECTEDAPPPLATFORMTYPE
        case "windows":
            result = WINDOWS_DETECTEDAPPPLATFORMTYPE
        case "windowsMobile":
            result = WINDOWSMOBILE_DETECTEDAPPPLATFORMTYPE
        case "windowsHolographic":
            result = WINDOWSHOLOGRAPHIC_DETECTEDAPPPLATFORMTYPE
        case "ios":
            result = IOS_DETECTEDAPPPLATFORMTYPE
        case "macOS":
            result = MACOS_DETECTEDAPPPLATFORMTYPE
        case "chromeOS":
            result = CHROMEOS_DETECTEDAPPPLATFORMTYPE
        case "androidOSP":
            result = ANDROIDOSP_DETECTEDAPPPLATFORMTYPE
        case "androidDeviceAdministrator":
            result = ANDROIDDEVICEADMINISTRATOR_DETECTEDAPPPLATFORMTYPE
        case "androidWorkProfile":
            result = ANDROIDWORKPROFILE_DETECTEDAPPPLATFORMTYPE
        case "androidDedicatedAndFullyManaged":
            result = ANDROIDDEDICATEDANDFULLYMANAGED_DETECTEDAPPPLATFORMTYPE
        default:
            return 0, errors.New("Unknown DetectedAppPlatformType value: " + v)
    }
    return &result, nil
}
func SerializeDetectedAppPlatformType(values []DetectedAppPlatformType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
