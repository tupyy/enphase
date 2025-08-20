package entity

import (
	"encoding/json"

	v1 "github.com/tupyy/enphase/api/v1"
)

// EnergyData represents comprehensive energy data domain entity
type EnergyData struct {
	// Production energy data
	Production *ProductionEnergy
	// Consumption energy data
	Consumption *ConsumptionEnergy
	// Raw holds the original JSON response
	Raw json.RawMessage
}

// ProductionEnergy represents production energy measurements
type ProductionEnergy struct {
	// PCU microinverters measurements
	PCU *EnergyMeasurement
	// RGM revenue grade meters (usually not applicable)
	RGM *EnergyMeasurement
	// EIM Envoy internal production meter
	EIM *EnergyMeasurement
}

// ConsumptionEnergy represents consumption energy measurements
type ConsumptionEnergy struct {
	// EIM Envoy internal consumption meter
	EIM *EnergyMeasurement
}

// EnergyMeasurement represents energy measurement data
type EnergyMeasurement struct {
	EnergyToday     int // Total energy today (Wh)
	EnergySevenDays int // Total energy last 7 days (Wh)
	EnergyLifetime  int // Total energy lifetime (Wh)
	PowerNow        int // Current active power (Watts)
}

// NewEnergyData creates a new EnergyData entity from API v1 model
func NewEnergyData(model *v1.EnergyData, rawJSON []byte) *EnergyData {
	entity := &EnergyData{
		Raw: json.RawMessage(rawJSON),
	}

	if model.Production != nil {
		entity.Production = &ProductionEnergy{}

		if model.Production.Pcu != nil {
			entity.Production.PCU = &EnergyMeasurement{}
			if model.Production.Pcu.WattHoursToday != nil {
				entity.Production.PCU.EnergyToday = *model.Production.Pcu.WattHoursToday
			}
			if model.Production.Pcu.WattHoursSevenDays != nil {
				entity.Production.PCU.EnergySevenDays = *model.Production.Pcu.WattHoursSevenDays
			}
			if model.Production.Pcu.WattHoursLifetime != nil {
				entity.Production.PCU.EnergyLifetime = *model.Production.Pcu.WattHoursLifetime
			}
			if model.Production.Pcu.WattsNow != nil {
				entity.Production.PCU.PowerNow = *model.Production.Pcu.WattsNow
			}
		}

		if model.Production.Rgm != nil {
			entity.Production.RGM = &EnergyMeasurement{}
			if model.Production.Rgm.WattHoursToday != nil {
				entity.Production.RGM.EnergyToday = *model.Production.Rgm.WattHoursToday
			}
			if model.Production.Rgm.WattHoursSevenDays != nil {
				entity.Production.RGM.EnergySevenDays = *model.Production.Rgm.WattHoursSevenDays
			}
			if model.Production.Rgm.WattHoursLifetime != nil {
				entity.Production.RGM.EnergyLifetime = *model.Production.Rgm.WattHoursLifetime
			}
			if model.Production.Rgm.WattsNow != nil {
				entity.Production.RGM.PowerNow = *model.Production.Rgm.WattsNow
			}
		}

		if model.Production.Eim != nil {
			entity.Production.EIM = &EnergyMeasurement{}
			if model.Production.Eim.WattHoursToday != nil {
				entity.Production.EIM.EnergyToday = *model.Production.Eim.WattHoursToday
			}
			if model.Production.Eim.WattHoursSevenDays != nil {
				entity.Production.EIM.EnergySevenDays = *model.Production.Eim.WattHoursSevenDays
			}
			if model.Production.Eim.WattHoursLifetime != nil {
				entity.Production.EIM.EnergyLifetime = *model.Production.Eim.WattHoursLifetime
			}
			if model.Production.Eim.WattsNow != nil {
				entity.Production.EIM.PowerNow = *model.Production.Eim.WattsNow
			}
		}
	}

	if model.Consumption != nil {
		entity.Consumption = &ConsumptionEnergy{}

		if model.Consumption.Eim != nil {
			entity.Consumption.EIM = &EnergyMeasurement{}
			if model.Consumption.Eim.WattHoursToday != nil {
				entity.Consumption.EIM.EnergyToday = *model.Consumption.Eim.WattHoursToday
			}
			if model.Consumption.Eim.WattHoursSevenDays != nil {
				entity.Consumption.EIM.EnergySevenDays = *model.Consumption.Eim.WattHoursSevenDays
			}
			if model.Consumption.Eim.WattHoursLifetime != nil {
				entity.Consumption.EIM.EnergyLifetime = *model.Consumption.Eim.WattHoursLifetime
			}
			if model.Consumption.Eim.WattsNow != nil {
				entity.Consumption.EIM.PowerNow = *model.Consumption.Eim.WattsNow
			}
		}
	}

	return entity
}
