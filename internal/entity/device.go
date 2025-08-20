package entity

import (
	"encoding/json"

	v1 "github.com/tupyy/enphase/api/v1"
)

// DeviceList represents ensemble device list domain entity
type DeviceList struct {
	// USB connection information
	USB *USBInfo
	// Devices list of provisioned devices
	Devices []Device
	// Raw holds the original JSON response
	Raw json.RawMessage
}

// USBInfo represents USB connection information
type USBInfo struct {
	CK2Bridge string // Connection status
	AutoScan  bool   // Auto scan enabled
}

// Device represents a provisioned device
type Device struct {
	SerialNumber     string // Device serial number
	DeviceType       int    // Device type (0=Unknown, 13=Storage, 14=Microinverters)
	DeviceTypeString string // Human readable device type
	ComInterface     int    // Communication interface number
	ComInterfaceStr  string // Communication interface string
	Status           string // Device connection status
	DeviceInfo       *DeviceInfo
}

// DeviceInfo represents device-specific information
type DeviceInfo struct {
	Capacity int // Device capacity in watt-hours
	DERIndex int // Connected phase (1=Phase1, 2=Phase2, 3=Phase3)
}

// NewDeviceList creates a new DeviceList entity from API v1 model
func NewDeviceList(model *v1.DeviceList, rawJSON []byte) *DeviceList {
	entity := &DeviceList{
		Raw: json.RawMessage(rawJSON),
	}

	if model.Usb != nil {
		entity.USB = &USBInfo{}
		if model.Usb.AutoScan != nil {
			entity.USB.AutoScan = *model.Usb.AutoScan == "true"
		}
		if model.Usb.Ck2Bridge != nil {
			entity.USB.CK2Bridge = *model.Usb.Ck2Bridge
		}
	}

	if model.Devices != nil {
		entity.Devices = make([]Device, len(*model.Devices))
		for i, device := range *model.Devices {
			entityDevice := Device{}

			if device.SerialNumber != nil {
				entityDevice.SerialNumber = *device.SerialNumber
			}
			if device.DeviceType != nil {
				entityDevice.DeviceType = *device.DeviceType
				entityDevice.DeviceTypeString = getDeviceTypeString(*device.DeviceType)
			}
			if device.ComInterface != nil {
				entityDevice.ComInterface = *device.ComInterface
			}
			if device.ComInterfaceStr != nil {
				entityDevice.ComInterfaceStr = *device.ComInterfaceStr
			}
			if device.Status != nil {
				entityDevice.Status = *device.Status
			}

			if device.DevInfo != nil {
				entityDevice.DeviceInfo = &DeviceInfo{}
				if device.DevInfo.Capacity != nil {
					entityDevice.DeviceInfo.Capacity = *device.DevInfo.Capacity
				}
				if device.DevInfo.DERIndex != nil {
					entityDevice.DeviceInfo.DERIndex = *device.DevInfo.DERIndex
				}
			}

			entity.Devices[i] = entityDevice
		}
	}

	return entity
}

// getDeviceTypeString returns human readable device type
func getDeviceTypeString(deviceType int) string {
	switch deviceType {
	case 0:
		return "Unknown"
	case 13:
		return "Storage"
	case 14:
		return "Microinverters"
	default:
		return "Unknown"
	}
}
