package entity

import (
	"encoding/json"
	"time"

	v1 "github.com/tupyy/enphase/api/v1"
)

// ProductionData represents production meter data domain entity
type ProductionData struct {
	EnergyToday     int // Total energy produced today (Wh)
	EnergySevenDays int // Total energy produced in last 7 days (Wh)
	EnergyLifetime  int // Total energy produced lifetime (Wh)
	PowerNow        int // Current active power production (Watts)
	// Raw holds the original JSON response
	Raw json.RawMessage
}

// InverterProduction represents inverter production data domain entity
type InverterProduction struct {
	Inverters []Inverter
	// Raw holds the original JSON response
	Raw json.RawMessage
}

// Inverter represents individual microinverter data
type Inverter struct {
	SerialNumber     string    // Serial number of microinverter
	LastReportDate   time.Time // Last report time
	DeviceType       int       // Device type
	LastReportWatts  int       // Last reported active power (Watts)
	MaxReportWatts   int       // Maximum ever reported power (Watts)
	IsOnline         bool      // Whether inverter is currently reporting
	EfficiencyRating float32   // Current vs max power ratio
}

// NewProductionData creates a new ProductionData entity from API v1 model
func NewProductionData(model *v1.ProductionData, rawJSON []byte) *ProductionData {
	entity := &ProductionData{
		Raw: json.RawMessage(rawJSON),
	}

	if model.WattHoursToday != nil {
		entity.EnergyToday = *model.WattHoursToday
	}
	if model.WattHoursSevenDays != nil {
		entity.EnergySevenDays = *model.WattHoursSevenDays
	}
	if model.WattHoursLifetime != nil {
		entity.EnergyLifetime = *model.WattHoursLifetime
	}
	if model.WattsNow != nil {
		entity.PowerNow = *model.WattsNow
	}

	return entity
}

// NewInverterProduction creates a new InverterProduction entity from API v1 model
func NewInverterProduction(model *v1.InverterProduction, rawJSON []byte) *InverterProduction {
	entity := &InverterProduction{
		Raw: json.RawMessage(rawJSON),
	}

	if model != nil {
		entity.Inverters = make([]Inverter, len(*model))
		for i, inverter := range *model {
			entityInverter := Inverter{}

			if inverter.SerialNumber != nil {
				entityInverter.SerialNumber = *inverter.SerialNumber
			}
			if inverter.LastReportDate != nil {
				entityInverter.LastReportDate = time.Unix(*inverter.LastReportDate, 0)
				// Consider inverter online if last report was within 10 minutes
				entityInverter.IsOnline = time.Since(entityInverter.LastReportDate) < 10*time.Minute
			}
			if inverter.DevType != nil {
				entityInverter.DeviceType = *inverter.DevType
			}
			if inverter.LastReportWatts != nil {
				entityInverter.LastReportWatts = *inverter.LastReportWatts
			}
			if inverter.MaxReportWatts != nil {
				entityInverter.MaxReportWatts = *inverter.MaxReportWatts
				// Calculate efficiency rating
				if entityInverter.MaxReportWatts > 0 {
					entityInverter.EfficiencyRating = float32(entityInverter.LastReportWatts) / float32(entityInverter.MaxReportWatts)
				}
			}

			entity.Inverters[i] = entityInverter
		}
	}

	return entity
}

// TotalCurrentPower returns the sum of current power from all inverters
func (ip *InverterProduction) TotalCurrentPower() int {
	var total int
	for _, inverter := range ip.Inverters {
		total += inverter.LastReportWatts
	}
	return total
}

// TotalMaxPower returns the sum of max power from all inverters
func (ip *InverterProduction) TotalMaxPower() int {
	var total int
	for _, inverter := range ip.Inverters {
		total += inverter.MaxReportWatts
	}
	return total
}

// OnlineInverters returns the count of online inverters
func (ip *InverterProduction) OnlineInverters() int {
	var count int
	for _, inverter := range ip.Inverters {
		if inverter.IsOnline {
			count++
		}
	}
	return count
}

// AverageEfficiency returns the average efficiency across all inverters
func (ip *InverterProduction) AverageEfficiency() float32 {
	if len(ip.Inverters) == 0 {
		return 0
	}

	var total float32
	for _, inverter := range ip.Inverters {
		total += inverter.EfficiencyRating
	}
	return total / float32(len(ip.Inverters))
}
