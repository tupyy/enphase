package presenter

import (
	"encoding/json"
	"fmt"

	"github.com/tupyy/enphase/internal/entity"
)

// ProductionPresenter handles presentation of production data
type ProductionPresenter struct{}

// NewProductionPresenter creates a new production presenter
func NewProductionPresenter() *ProductionPresenter {
	return &ProductionPresenter{}
}

// Summary presents production data in a summary format
func (p *ProductionPresenter) Summary(data *entity.ProductionData) error {
	// TODO: Implement summary presentation
	fmt.Println("Summary presentation not yet implemented")
	return nil
}

// Raw presents production data in raw JSON format
func (p *ProductionPresenter) Raw(data *entity.ProductionData) error {
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

// InverterPresenter handles presentation of inverter production data
type InverterPresenter struct{}

// NewInverterPresenter creates a new inverter presenter
func NewInverterPresenter() *InverterPresenter {
	return &InverterPresenter{}
}

// Summary presents inverter production data in a summary format
func (p *InverterPresenter) Summary(data *entity.InverterProduction) error {
	// TODO: Implement summary presentation
	fmt.Println("Summary presentation not yet implemented")
	return nil
}

// Raw presents inverter production data in raw JSON format
func (p *InverterPresenter) Raw(data *entity.InverterProduction) error {
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
