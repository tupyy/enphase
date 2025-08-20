package presenter

import (
	"encoding/json"
	"fmt"

	"github.com/tupyy/enphase/internal/entity"
)

// ConsumptionPresenter handles presentation of consumption data
type ConsumptionPresenter struct{}

// NewConsumptionPresenter creates a new consumption presenter
func NewConsumptionPresenter() *ConsumptionPresenter {
	return &ConsumptionPresenter{}
}

// Summary presents consumption data in a summary format
func (p *ConsumptionPresenter) Summary(data *entity.ConsumptionData) error {
	// TODO: Implement summary presentation
	fmt.Println("Summary presentation not yet implemented")
	return nil
}

// Raw presents consumption data in raw JSON format
func (p *ConsumptionPresenter) Raw(data *entity.ConsumptionData) error {
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
