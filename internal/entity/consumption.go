package entity

import (
	"encoding/json"
	"time"

	v1 "github.com/tupyy/enphase/api/v1"
)

// ConsumptionData represents power consumption information domain entity
type ConsumptionData struct {
	// CreatedAt timestamp of measurement
	CreatedAt time.Time
	// ReportType type of consumption report
	ReportType string
	// Cumulative results of all phases
	Cumulative *ConsumptionCumulative
	// Lines per-phase measurements
	Lines []ConsumptionLine
	// Raw holds the original JSON response
	Raw json.RawMessage
}

// ConsumptionCumulative represents cumulative consumption data
type ConsumptionCumulative struct {
	InstantaneousActivePower float32 // Watts
	ActivePower              float32 // Watts
	ApparentPower            float32 // VA
	ReactivePower            float32 // VAr
	EnergyDeliveredCum       float32 // Wh
	EnergyReceivedCum        float32 // Wh
	ReactiveEnergyLagCum     float32 // VArh
	ReactiveEnergyLeadCum    float32 // VArh
	ApparentEnergyCum        float32 // VAh
	RMSVoltage               float32 // Volts
	RMSCurrent               float32 // Amperes
	PowerFactor              float32
	Frequency                float32 // Hz
}

// ConsumptionLine represents per-phase consumption measurements
type ConsumptionLine struct {
	InstantaneousActivePower float32 // Watts
	ActivePower              float32 // Watts
	ApparentPower            float32 // VA
	ReactivePower            float32 // VAr
	EnergyDeliveredCum       float32 // Wh
	EnergyReceivedCum        float32 // Wh
	ReactiveEnergyLagCum     float32 // VArh
	ReactiveEnergyLeadCum    float32 // VArh
	ApparentEnergyCum        float32 // VAh
	RMSVoltage               float32 // Volts
	RMSCurrent               float32 // Amperes
	PowerFactor              float32
	Frequency                float32 // Hz
}

// NewConsumptionData creates a new ConsumptionData entity from API v1 model
func NewConsumptionData(model *v1.ConsumptionData, rawJSON []byte) *ConsumptionData {
	entity := &ConsumptionData{
		Raw: json.RawMessage(rawJSON),
	}

	if model.CreatedAt != nil {
		entity.CreatedAt = time.Unix(*model.CreatedAt, 0)
	}

	if model.ReportType != nil {
		entity.ReportType = string(*model.ReportType)
	}

	if model.Cumulative != nil {
		entity.Cumulative = &ConsumptionCumulative{}
		if model.Cumulative.CurrW != nil {
			entity.Cumulative.InstantaneousActivePower = *model.Cumulative.CurrW
		}
		if model.Cumulative.ActPower != nil {
			entity.Cumulative.ActivePower = *model.Cumulative.ActPower
		}
		if model.Cumulative.ApprntPwr != nil {
			entity.Cumulative.ApparentPower = *model.Cumulative.ApprntPwr
		}
		if model.Cumulative.ReactPwr != nil {
			entity.Cumulative.ReactivePower = *model.Cumulative.ReactPwr
		}
		if model.Cumulative.WhDlvdCum != nil {
			entity.Cumulative.EnergyDeliveredCum = *model.Cumulative.WhDlvdCum
		}
		if model.Cumulative.WhRcvdCum != nil {
			entity.Cumulative.EnergyReceivedCum = *model.Cumulative.WhRcvdCum
		}
		if model.Cumulative.VarhLagCum != nil {
			entity.Cumulative.ReactiveEnergyLagCum = *model.Cumulative.VarhLagCum
		}
		if model.Cumulative.VarhLeadCum != nil {
			entity.Cumulative.ReactiveEnergyLeadCum = *model.Cumulative.VarhLeadCum
		}
		if model.Cumulative.VahCum != nil {
			entity.Cumulative.ApparentEnergyCum = *model.Cumulative.VahCum
		}
		if model.Cumulative.RmsVoltage != nil {
			entity.Cumulative.RMSVoltage = *model.Cumulative.RmsVoltage
		}
		if model.Cumulative.RmsCurrent != nil {
			entity.Cumulative.RMSCurrent = *model.Cumulative.RmsCurrent
		}
		if model.Cumulative.PwrFactor != nil {
			entity.Cumulative.PowerFactor = *model.Cumulative.PwrFactor
		}
		if model.Cumulative.FreqHz != nil {
			entity.Cumulative.Frequency = *model.Cumulative.FreqHz
		}
	}

	if model.Lines != nil {
		entity.Lines = make([]ConsumptionLine, len(*model.Lines))
		for i, line := range *model.Lines {
			entityLine := ConsumptionLine{}
			if line.CurrW != nil {
				entityLine.InstantaneousActivePower = *line.CurrW
			}
			if line.ActPower != nil {
				entityLine.ActivePower = *line.ActPower
			}
			if line.ApprntPwr != nil {
				entityLine.ApparentPower = *line.ApprntPwr
			}
			if line.ReactPwr != nil {
				entityLine.ReactivePower = *line.ReactPwr
			}
			if line.WhDlvdCum != nil {
				entityLine.EnergyDeliveredCum = *line.WhDlvdCum
			}
			if line.WhRcvdCum != nil {
				entityLine.EnergyReceivedCum = *line.WhRcvdCum
			}
			if line.VarhLagCum != nil {
				entityLine.ReactiveEnergyLagCum = *line.VarhLagCum
			}
			if line.VarhLeadCum != nil {
				entityLine.ReactiveEnergyLeadCum = *line.VarhLeadCum
			}
			if line.VahCum != nil {
				entityLine.ApparentEnergyCum = *line.VahCum
			}
			if line.RmsVoltage != nil {
				entityLine.RMSVoltage = *line.RmsVoltage
			}
			if line.RmsCurrent != nil {
				entityLine.RMSCurrent = *line.RmsCurrent
			}
			if line.PwrFactor != nil {
				entityLine.PowerFactor = *line.PwrFactor
			}
			if line.FreqHz != nil {
				entityLine.Frequency = *line.FreqHz
			}
			entity.Lines[i] = entityLine
		}
	}

	return entity
}
