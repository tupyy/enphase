package entity

import (
	"encoding/json"

	v1 "github.com/tupyy/enphase/api/v1"
)

// GridReading represents grid connection readings domain entity
type GridReading struct {
	// Channels per-phase grid readings
	Channels []GridChannel
	// Raw holds the original JSON response
	Raw json.RawMessage
}

// GridChannel represents per-phase grid readings
type GridChannel struct {
	Phase         string  // Phase identifier (L1, L2, L3)
	ActivePower   float32 // Instantaneous active power (Watts, negative=export)
	ReactivePower float32 // Instantaneous reactive power (VAr)
	Voltage       float32 // Measured voltage (Volts)
	Current       float32 // Measured current (Amperes)
	Frequency     float32 // Measured frequency (Hz)
}

// NewGridReading creates a new GridReading entity from API v1 model
func NewGridReading(model *v1.GridReading, rawJSON []byte) *GridReading {
	entity := &GridReading{
		Raw: json.RawMessage(rawJSON),
	}

	if model != nil && len(*model) > 0 {
		// GridReading is an array, take the first element
		firstElement := (*model)[0]

		if firstElement.Channels != nil {
			entity.Channels = make([]GridChannel, len(*firstElement.Channels))
			for i, channel := range *firstElement.Channels {
				entityChannel := GridChannel{}

				if channel.Phase != nil {
					entityChannel.Phase = string(*channel.Phase)
				}
				if channel.ActivePower != nil {
					entityChannel.ActivePower = *channel.ActivePower
				}
				if channel.ReactivePower != nil {
					entityChannel.ReactivePower = *channel.ReactivePower
				}
				if channel.Voltage != nil {
					entityChannel.Voltage = *channel.Voltage
				}
				if channel.Current != nil {
					entityChannel.Current = *channel.Current
				}
				if channel.Freq != nil {
					entityChannel.Frequency = *channel.Freq
				}

				entity.Channels[i] = entityChannel
			}
		}
	}

	return entity
}

// IsExporting returns true if the grid is exporting power (negative values)
func (g *GridReading) IsExporting() bool {
	for _, channel := range g.Channels {
		if channel.ActivePower < 0 {
			return true
		}
	}
	return false
}

// TotalActivePower returns the sum of active power across all phases
func (g *GridReading) TotalActivePower() float32 {
	var total float32
	for _, channel := range g.Channels {
		total += channel.ActivePower
	}
	return total
}

// TotalReactivePower returns the sum of reactive power across all phases
func (g *GridReading) TotalReactivePower() float32 {
	var total float32
	for _, channel := range g.Channels {
		total += channel.ReactivePower
	}
	return total
}
