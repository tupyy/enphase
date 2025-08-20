package presenter

import (
	"encoding/json"
	"fmt"

	"github.com/tupyy/enphase/internal/entity"
)

// LiveDataPresenter handles presentation of live data status
type LiveDataPresenter struct{}

// NewLiveDataPresenter creates a new live data presenter
func NewLiveDataPresenter() *LiveDataPresenter {
	return &LiveDataPresenter{}
}

// Summary presents live data status in a summary format
func (p *LiveDataPresenter) Summary(data *entity.LiveDataStatus) error {
	// TODO: Implement summary presentation
	fmt.Println("Summary presentation not yet implemented")
	return nil
}

// Raw presents live data status in raw JSON format
func (p *LiveDataPresenter) Raw(data *entity.LiveDataStatus) error {
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
