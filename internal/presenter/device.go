package presenter

import (
	"encoding/json"
	"fmt"

	"github.com/tupyy/enphase/internal/entity"
)

// DevicePresenter handles presentation of device list data
type DevicePresenter struct{}

// NewDevicePresenter creates a new device presenter
func NewDevicePresenter() *DevicePresenter {
	return &DevicePresenter{}
}

// Summary presents device list data in a summary format
func (p *DevicePresenter) Summary(data *entity.DeviceList) error {
	// TODO: Implement summary presentation
	fmt.Println("Summary presentation not yet implemented")
	return nil
}

// Raw presents device list data in raw JSON format
func (p *DevicePresenter) Raw(data *entity.DeviceList) error {
	if data.Raw == nil {
		return fmt.Errorf("no raw data available")
	}

	var jsonData interface{}
	if err := json.Unmarshal(data.Raw, &jsonData); err != nil {
		return fmt.Errorf("failed to unmarshal raw data: %w", err)
	}

	prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON: %w", err)
	}

	fmt.Println(string(prettyJSON))
	return nil
}
