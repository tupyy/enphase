package presenter

import (
	"encoding/json"
	"fmt"

	"github.com/tupyy/enphase/internal/entity"
)

// MeterPresenter handles presentation of meter details
type MeterPresenter struct{}

// NewMeterPresenter creates a new meter presenter
func NewMeterPresenter() *MeterPresenter {
	return &MeterPresenter{}
}

// Summary presents meter details in a summary format
func (p *MeterPresenter) Summary(data *entity.MeterDetails) error {
	// TODO: Implement summary presentation
	fmt.Println("Summary presentation not yet implemented")
	return nil
}

// Raw presents meter details in raw JSON format
func (p *MeterPresenter) Raw(data *entity.MeterDetails) error {
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

// MeterReadingsPresenter handles presentation of meter readings
type MeterReadingsPresenter struct{}

// NewMeterReadingsPresenter creates a new meter readings presenter
func NewMeterReadingsPresenter() *MeterReadingsPresenter {
	return &MeterReadingsPresenter{}
}

// Summary presents meter readings in a summary format
func (p *MeterReadingsPresenter) Summary(data *entity.MeterReadings) error {
	// TODO: Implement summary presentation
	fmt.Println("Summary presentation not yet implemented")
	return nil
}

// Raw presents meter readings in raw JSON format
func (p *MeterReadingsPresenter) Raw(data *entity.MeterReadings) error {
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
