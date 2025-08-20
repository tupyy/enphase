package presenter

import (
	"encoding/json"
	"fmt"

	"github.com/tupyy/enphase/internal/entity"
)

// GridPresenter handles presentation of grid reading data
type GridPresenter struct{}

// NewGridPresenter creates a new grid presenter
func NewGridPresenter() *GridPresenter {
	return &GridPresenter{}
}

// Summary presents grid reading data in a summary format
func (p *GridPresenter) Summary(data *entity.GridReading) error {
	// TODO: Implement summary presentation
	fmt.Println("Summary presentation not yet implemented")
	return nil
}

// Raw presents grid reading data in raw JSON format
func (p *GridPresenter) Raw(data *entity.GridReading) error {
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
