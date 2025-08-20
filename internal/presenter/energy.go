package presenter

import (
	"encoding/json"
	"fmt"

	"github.com/tupyy/enphase/internal/entity"
)

// EnergyPresenter handles presentation of energy data
type EnergyPresenter struct{}

// NewEnergyPresenter creates a new energy presenter
func NewEnergyPresenter() *EnergyPresenter {
	return &EnergyPresenter{}
}

// Summary presents energy data in a summary format
func (p *EnergyPresenter) Summary(data *entity.EnergyData) error {
	// TODO: Implement summary presentation
	fmt.Println("Summary presentation not yet implemented")
	return nil
}

// Raw presents energy data in raw JSON format
func (p *EnergyPresenter) Raw(data *entity.EnergyData) error {
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
