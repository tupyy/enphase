package entity

import (
	"encoding/json"
	"time"

	v1 "github.com/tupyy/enphase/api/v1"
)

// ConsumptionData represents power consumption information domain entity
type ConsumptionData struct {
	// Items holds all consumption data entries from the array
	Items []ConsumptionItem
	// Raw holds the original JSON response
	Raw json.RawMessage
}

// ConsumptionItem represents a single consumption data entry
type ConsumptionItem struct {
	// CreatedAt timestamp of measurement
	CreatedAt time.Time
	// ReportType type of consumption report
	ReportType string
	// Cumulative results of all phases
	Cumulative *ConsumptionCumulative
	// Lines per-phase measurements
	Lines []ConsumptionLine
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

// NewConsumptionData creates a new ConsumptionData entity from API v1 model array
func NewConsumptionData(models []v1.ConsumptionData, rawJSON []byte) *ConsumptionData {
	entity := &ConsumptionData{
		Raw:   json.RawMessage(rawJSON),
		Items: make([]ConsumptionItem, len(models)),
	}

	for idx, model := range models {
		item := ConsumptionItem{}

		if model.CreatedAt != nil {
			item.CreatedAt = time.Unix(*model.CreatedAt, 0)
		}

		if model.ReportType != nil {
			item.ReportType = string(*model.ReportType)
		}

		if model.Cumulative != nil {
			item.Cumulative = &ConsumptionCumulative{}
			if model.Cumulative.CurrW != nil {
				item.Cumulative.InstantaneousActivePower = *model.Cumulative.CurrW
			}
			if model.Cumulative.ActPower != nil {
				item.Cumulative.ActivePower = *model.Cumulative.ActPower
			}
			if model.Cumulative.ApprntPwr != nil {
				item.Cumulative.ApparentPower = *model.Cumulative.ApprntPwr
			}
			if model.Cumulative.ReactPwr != nil {
				item.Cumulative.ReactivePower = *model.Cumulative.ReactPwr
			}
			if model.Cumulative.WhDlvdCum != nil {
				item.Cumulative.EnergyDeliveredCum = *model.Cumulative.WhDlvdCum
			}
			if model.Cumulative.WhRcvdCum != nil {
				item.Cumulative.EnergyReceivedCum = *model.Cumulative.WhRcvdCum
			}
			if model.Cumulative.VarhLagCum != nil {
				item.Cumulative.ReactiveEnergyLagCum = *model.Cumulative.VarhLagCum
			}
			if model.Cumulative.VarhLeadCum != nil {
				item.Cumulative.ReactiveEnergyLeadCum = *model.Cumulative.VarhLeadCum
			}
			if model.Cumulative.VahCum != nil {
				item.Cumulative.ApparentEnergyCum = *model.Cumulative.VahCum
			}
			if model.Cumulative.RmsVoltage != nil {
				item.Cumulative.RMSVoltage = *model.Cumulative.RmsVoltage
			}
			if model.Cumulative.RmsCurrent != nil {
				item.Cumulative.RMSCurrent = *model.Cumulative.RmsCurrent
			}
			if model.Cumulative.PwrFactor != nil {
				item.Cumulative.PowerFactor = *model.Cumulative.PwrFactor
			}
			if model.Cumulative.FreqHz != nil {
				item.Cumulative.Frequency = *model.Cumulative.FreqHz
			}
		}

		if model.Lines != nil {
			item.Lines = make([]ConsumptionLine, len(*model.Lines))
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
				item.Lines[i] = entityLine
			}
		}

		entity.Items[idx] = item
	}

	return entity
}
