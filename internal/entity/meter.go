package entity

import (
	"encoding/json"
	"time"

	v1 "github.com/tupyy/enphase/api/v1"
)

// MeterDetails represents meter configuration domain entity
type MeterDetails struct {
	Meters []Meter
	// Raw holds the original JSON response
	Raw json.RawMessage
}

// Meter represents individual meter configuration
type Meter struct {
	ID              int64  // Meter ID
	State           string // Meter state (enabled/disabled)
	MeasurementType string // Type of measurement
	PhaseMode       string // Number of phases the meter can monitor
	PhaseCount      int    // Number of phases meter is monitoring
	MeteringStatus  string // Metering status
	StatusFlags     []string
	IsEnabled       bool // Convenience field
	IsProduction    bool // True if this is a production meter
	IsConsumption   bool // True if this is a consumption meter
	IsStorage       bool // True if this is a storage meter
}

// MeterReadings represents detailed meter readings domain entity
type MeterReadings struct {
	Readings []MeterReading
	// Raw holds the original JSON response
	Raw json.RawMessage
}

// MeterReading represents individual meter reading
type MeterReading struct {
	ID                    int64          // Meter ID
	Timestamp             time.Time      // Measurement timestamp
	ActiveEnergyDelivered float32        // Overall active energy delivered (Wh)
	ActiveEnergyReceived  float32        // Overall active energy received (Wh)
	ApparentEnergy        float32        // Overall apparent energy (VAh)
	ReactiveEnergyLagging float32        // Overall lagging reactive energy (VArh)
	ReactiveEnergyLeading float32        // Overall leading reactive energy (VArh)
	InstantaneousDemand   float32        // Instantaneous active power demand (Watts)
	ActivePower           float32        // Instantaneous active power (Watts)
	ApparentPower         float32        // Instantaneous apparent power (VA)
	ReactivePower         float32        // Instantaneous reactive power (VAr)
	PowerFactor           float32        // Measured power factor
	Voltage               float32        // Measured voltage (Volts)
	Current               float32        // Measured current (Amperes)
	Frequency             float32        // Measured frequency (Hz)
	Channels              []MeterChannel // Per-phase channel measurements
	MeterType             string         // Derived meter type for convenience
	IsExporting           bool           // True if exporting power
}

// MeterChannel represents per-phase meter channel reading
type MeterChannel struct {
	ID                    int64     // Channel ID
	Timestamp             time.Time // Channel measurement timestamp
	ActiveEnergyDelivered float32   // Active energy delivered in this channel (Wh)
	ActiveEnergyReceived  float32   // Active energy received in this channel (Wh)
	ApparentEnergy        float32   // Apparent energy in this channel (VAh)
	ReactiveEnergyLagging float32   // Lagging reactive energy in this channel (VArh)
	ReactiveEnergyLeading float32   // Leading reactive energy in this channel (VArh)
	InstantaneousDemand   float32   // Instantaneous active power demand in this channel (Watts)
	ActivePower           float32   // Instantaneous active power in this channel (Watts)
	ApparentPower         float32   // Instantaneous apparent power in this channel (VA)
	ReactivePower         float32   // Instantaneous reactive power in this channel (VAr)
	PowerFactor           float32   // Measured power factor in this channel
	Voltage               float32   // Measured voltage in this channel (Volts)
	Current               float32   // Measured current in this channel (Amperes)
	Frequency             float32   // Measured frequency in this channel (Hz)
	Phase                 string    // Phase identifier (derived from channel number)
}

// NewMeterDetails creates a new MeterDetails entity from API v1 model
func NewMeterDetails(model *v1.MeterDetails, rawJSON []byte) *MeterDetails {
	entity := &MeterDetails{
		Raw: json.RawMessage(rawJSON),
	}

	if model != nil {
		entity.Meters = make([]Meter, len(*model))
		for i, meter := range *model {
			entityMeter := Meter{}

			if meter.Eid != nil {
				entityMeter.ID = *meter.Eid
			}
			if meter.State != nil {
				entityMeter.State = *meter.State
				entityMeter.IsEnabled = *meter.State == "enabled"
			}
			if meter.MeasurementType != nil {
				entityMeter.MeasurementType = string(*meter.MeasurementType)
				entityMeter.IsProduction = *meter.MeasurementType == v1.MeterDetailsMeasurementTypeProduction
				entityMeter.IsConsumption = *meter.MeasurementType == v1.MeterDetailsMeasurementTypeNetConsumption
				entityMeter.IsStorage = *meter.MeasurementType == v1.MeterDetailsMeasurementTypeStorage
			}
			if meter.PhaseMode != nil {
				entityMeter.PhaseMode = *meter.PhaseMode
			}
			if meter.PhaseCount != nil {
				entityMeter.PhaseCount = *meter.PhaseCount
			}
			if meter.MeteringStatus != nil {
				entityMeter.MeteringStatus = *meter.MeteringStatus
			}
			if meter.StatusFlags != nil {
				entityMeter.StatusFlags = *meter.StatusFlags
			}

			entity.Meters[i] = entityMeter
		}
	}

	return entity
}

// NewMeterReadings creates a new MeterReadings entity from API v1 model
func NewMeterReadings(model *v1.MeterReadings, rawJSON []byte) *MeterReadings {
	entity := &MeterReadings{
		Raw: json.RawMessage(rawJSON),
	}

	if model != nil {
		entity.Readings = make([]MeterReading, len(*model))
		for i, reading := range *model {
			entityReading := MeterReading{}

			if reading.Eid != nil {
				entityReading.ID = *reading.Eid
				entityReading.MeterType = deriveMeterType(*reading.Eid)
			}
			if reading.Timestamp != nil {
				entityReading.Timestamp = time.Unix(*reading.Timestamp, 0)
			}
			if reading.ActEnergyDlvd != nil {
				entityReading.ActiveEnergyDelivered = *reading.ActEnergyDlvd
			}
			if reading.ActEnergyRcvd != nil {
				entityReading.ActiveEnergyReceived = *reading.ActEnergyRcvd
			}
			if reading.ApparentEnergy != nil {
				entityReading.ApparentEnergy = *reading.ApparentEnergy
			}
			if reading.ReactEnergyLagg != nil {
				entityReading.ReactiveEnergyLagging = *reading.ReactEnergyLagg
			}
			if reading.ReactEnergyLead != nil {
				entityReading.ReactiveEnergyLeading = *reading.ReactEnergyLead
			}
			if reading.InstantaneousDemand != nil {
				entityReading.InstantaneousDemand = *reading.InstantaneousDemand
			}
			if reading.ActivePower != nil {
				entityReading.ActivePower = *reading.ActivePower
				entityReading.IsExporting = *reading.ActivePower < 0
			}
			if reading.ApparentPower != nil {
				entityReading.ApparentPower = *reading.ApparentPower
			}
			if reading.ReactivePower != nil {
				entityReading.ReactivePower = *reading.ReactivePower
			}
			if reading.PwrFactor != nil {
				entityReading.PowerFactor = *reading.PwrFactor
			}
			if reading.Voltage != nil {
				entityReading.Voltage = *reading.Voltage
			}
			if reading.Current != nil {
				entityReading.Current = *reading.Current
			}
			if reading.Freq != nil {
				entityReading.Frequency = *reading.Freq
			}

			// Process channels
			if reading.Channels != nil {
				entityReading.Channels = make([]MeterChannel, len(*reading.Channels))
				for j, channel := range *reading.Channels {
					entityChannel := MeterChannel{}

					if channel.Eid != nil {
						entityChannel.ID = *channel.Eid
						entityChannel.Phase = derivePhaseFromChannelID(*channel.Eid)
					}
					if channel.Timestamp != nil {
						entityChannel.Timestamp = time.Unix(*channel.Timestamp, 0)
					}
					if channel.ActEnergyDlvd != nil {
						entityChannel.ActiveEnergyDelivered = *channel.ActEnergyDlvd
					}
					if channel.ActEnergyRcvd != nil {
						entityChannel.ActiveEnergyReceived = *channel.ActEnergyRcvd
					}
					if channel.ApparentEnergy != nil {
						entityChannel.ApparentEnergy = *channel.ApparentEnergy
					}
					if channel.ReactEnergyLagg != nil {
						entityChannel.ReactiveEnergyLagging = *channel.ReactEnergyLagg
					}
					if channel.ReactEnergyLead != nil {
						entityChannel.ReactiveEnergyLeading = *channel.ReactEnergyLead
					}
					if channel.InstantaneousDemand != nil {
						entityChannel.InstantaneousDemand = *channel.InstantaneousDemand
					}
					if channel.ActivePower != nil {
						entityChannel.ActivePower = *channel.ActivePower
					}
					if channel.ApparentPower != nil {
						entityChannel.ApparentPower = *channel.ApparentPower
					}
					if channel.ReactivePower != nil {
						entityChannel.ReactivePower = *channel.ReactivePower
					}
					if channel.PwrFactor != nil {
						entityChannel.PowerFactor = *channel.PwrFactor
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

					entityReading.Channels[j] = entityChannel
				}
			}

			entity.Readings[i] = entityReading
		}
	}

	return entity
}

// deriveMeterType returns meter type based on meter ID patterns
func deriveMeterType(meterID int64) string {
	// These are typical patterns observed in Enphase systems
	switch {
	case meterID >= 704643328 && meterID < 704643584:
		return "Production"
	case meterID >= 704643584 && meterID < 704643840:
		return "Storage"
	case meterID >= 704643840:
		return "Consumption"
	default:
		return "Unknown"
	}
}

// derivePhaseFromChannelID returns phase identifier based on channel ID
func derivePhaseFromChannelID(channelID int64) string {
	// Channel IDs typically end in different patterns for different phases
	switch channelID % 10 {
	case 9:
		return "L1"
	case 0:
		return "L2"
	case 1:
		return "L3"
	default:
		return "Unknown"
	}
}

// GetProductionMeter returns the production meter if available
func (md *MeterDetails) GetProductionMeter() *Meter {
	for _, meter := range md.Meters {
		if meter.IsProduction {
			return &meter
		}
	}
	return nil
}

// GetConsumptionMeter returns the consumption meter if available
func (md *MeterDetails) GetConsumptionMeter() *Meter {
	for _, meter := range md.Meters {
		if meter.IsConsumption {
			return &meter
		}
	}
	return nil
}

// GetStorageMeter returns the storage meter if available
func (md *MeterDetails) GetStorageMeter() *Meter {
	for _, meter := range md.Meters {
		if meter.IsStorage {
			return &meter
		}
	}
	return nil
}
