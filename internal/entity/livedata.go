package entity

import (
	"encoding/json"
	"time"

	v1 "github.com/tupyy/enphase/api/v1"
)

// LiveDataStatus represents live system status domain entity
type LiveDataStatus struct {
	// Connection status information
	Connection *ConnectionStatus
	// Meters live data
	Meters *LiveMeters
	// Tasks information
	Tasks *TasksInfo
	// Counters system counters
	Counters map[string]int
	// Raw holds the original JSON response
	Raw json.RawMessage
}

// ConnectionStatus represents connection status information
type ConnectionStatus struct {
	MQTTState   string // MQTT broker status
	ProvState   string // Provisioning status
	AuthState   string // Authentication status
	StreamState string // Status of live data stream
	DebugState  string // Status of debug stream
	IsConnected bool   // Convenience field
}

// LiveMeters represents live meter data
type LiveMeters struct {
	LastUpdate         time.Time // Last meter data collection timestamp
	StateOfCharge      int       // State of charge of storage device (%)
	MainRelayState     int       // Main relay state
	GenRelayState      int       // Generator relay state
	BackupBatMode      int       // Backup battery mode
	BackupSOC          int       // Backup battery state of charge (%)
	IsSplitPhase       bool      // Split-phase grid type indicator
	PhaseCount         int       // Number of connected phases
	EnchargeAggSOC     int       // Aggregate SoC of Encharge (%)
	EnchargeAggEnergy  int       // Aggregate energy of Encharge (Wh)
	ACBatteryAggSOC    int       // Aggregate SoC of AC battery (%)
	ACBatteryAggEnergy int       // Aggregate energy of AC battery (Wh)

	// Meter readings
	PV        *PowerMeter // Photovoltaic meter readings
	Storage   *PowerMeter // Battery storage meter readings
	Grid      *PowerMeter // Grid meter readings
	Load      *PowerMeter // Load meter readings
	Generator *PowerMeter // Generator meter readings
}

// PowerMeter represents power meter readings
type PowerMeter struct {
	AggregateActivePower   int // Aggregate active power (milliwatts)
	AggregateApparentPower int // Aggregate apparent power (milli volt-amperes)
	ActivePowerPhaseA      int // Active power phase A (milliwatts)
	ActivePowerPhaseB      int // Active power phase B (milliwatts)
	ActivePowerPhaseC      int // Active power phase C (milliwatts)
	ApparentPowerPhaseA    int // Apparent power phase A (milli volt-amperes)
	ApparentPowerPhaseB    int // Apparent power phase B (milli volt-amperes)
	ApparentPowerPhaseC    int // Apparent power phase C (milli volt-amperes)

	// Convenience fields (converted to watts)
	ActivePowerWatts float32 // Active power in watts
	ApparentPowerVA  float32 // Apparent power in VA
	IsExporting      bool    // True if exporting power (negative)
}

// TasksInfo represents task information
type TasksInfo struct {
	TaskID    int       // ID of most recent task processed
	Timestamp time.Time // Unix timestamp of last task processed
}

// NewLiveDataStatus creates a new LiveDataStatus entity from API v1 model
func NewLiveDataStatus(model *v1.LiveDataStatus, rawJSON []byte) *LiveDataStatus {
	entity := &LiveDataStatus{
		Raw: json.RawMessage(rawJSON),
	}

	if model.Connection != nil {
		entity.Connection = &ConnectionStatus{}
		if model.Connection.MqttState != nil {
			entity.Connection.MQTTState = *model.Connection.MqttState
		}
		if model.Connection.ProvState != nil {
			entity.Connection.ProvState = *model.Connection.ProvState
		}
		if model.Connection.AuthState != nil {
			entity.Connection.AuthState = *model.Connection.AuthState
		}
		if model.Connection.ScStream != nil {
			entity.Connection.StreamState = *model.Connection.ScStream
		}
		if model.Connection.ScDebug != nil {
			entity.Connection.DebugState = *model.Connection.ScDebug
		}

		// Determine if connected
		entity.Connection.IsConnected = entity.Connection.MQTTState == "connected" &&
			entity.Connection.AuthState == "ok"
	}

	if model.Meters != nil {
		entity.Meters = &LiveMeters{}

		if model.Meters.LastUpdate != nil {
			entity.Meters.LastUpdate = time.Unix(*model.Meters.LastUpdate, 0)
		}
		if model.Meters.Soc != nil {
			entity.Meters.StateOfCharge = *model.Meters.Soc
		}
		if model.Meters.MainRelayState != nil {
			entity.Meters.MainRelayState = *model.Meters.MainRelayState
		}
		if model.Meters.GenRelayState != nil {
			entity.Meters.GenRelayState = *model.Meters.GenRelayState
		}
		if model.Meters.BackupBatMode != nil {
			entity.Meters.BackupBatMode = *model.Meters.BackupBatMode
		}
		if model.Meters.BackupSoc != nil {
			entity.Meters.BackupSOC = *model.Meters.BackupSoc
		}
		if model.Meters.IsSplitPhase != nil {
			entity.Meters.IsSplitPhase = *model.Meters.IsSplitPhase == 1
		}
		if model.Meters.PhaseCount != nil {
			entity.Meters.PhaseCount = *model.Meters.PhaseCount
		}
		if model.Meters.EncAggSoc != nil {
			entity.Meters.EnchargeAggSOC = *model.Meters.EncAggSoc
		}
		if model.Meters.EncAggEnergy != nil {
			entity.Meters.EnchargeAggEnergy = *model.Meters.EncAggEnergy
		}
		if model.Meters.AcbAggSoc != nil {
			entity.Meters.ACBatteryAggSOC = *model.Meters.AcbAggSoc
		}
		if model.Meters.AcbAggEnergy != nil {
			entity.Meters.ACBatteryAggEnergy = *model.Meters.AcbAggEnergy
		}

		// Convert meter readings
		if model.Meters.Pv != nil {
			entity.Meters.PV = convertPowerMeter(model.Meters.Pv)
		}
		if model.Meters.Storage != nil {
			entity.Meters.Storage = convertPowerMeter(model.Meters.Storage)
		}
		if model.Meters.Grid != nil {
			entity.Meters.Grid = convertPowerMeter(model.Meters.Grid)
		}
		if model.Meters.Load != nil {
			entity.Meters.Load = convertPowerMeter(model.Meters.Load)
		}
		if model.Meters.Generator != nil {
			entity.Meters.Generator = convertPowerMeter(model.Meters.Generator)
		}
	}

	if model.Tasks != nil {
		entity.Tasks = &TasksInfo{}
		if model.Tasks.TaskId != nil {
			entity.Tasks.TaskID = *model.Tasks.TaskId
		}
		if model.Tasks.Timestamp != nil {
			entity.Tasks.Timestamp = time.Unix(*model.Tasks.Timestamp, 0)
		}
	}

	if model.Counters != nil {
		entity.Counters = make(map[string]int)
		for key, value := range *model.Counters {
			entity.Counters[key] = value
		}
	}

	return entity
}

// convertPowerMeter converts API v1 power meter to entity power meter
func convertPowerMeter(apiMeter interface{}) *PowerMeter {
	// This is a generic converter since all meter types have the same structure
	// We'll use type assertion or reflection to handle the conversion

	meter := &PowerMeter{}

	// Use reflection or type assertion to extract values
	// For simplicity, we'll use a map-based approach
	jsonData, _ := json.Marshal(apiMeter)
	var meterMap map[string]*int
	json.Unmarshal(jsonData, &meterMap)

	if aggPMw, ok := meterMap["agg_p_mw"]; ok && aggPMw != nil {
		meter.AggregateActivePower = *aggPMw
		meter.ActivePowerWatts = float32(*aggPMw) / 1000.0 // Convert milliwatts to watts
		meter.IsExporting = *aggPMw < 0
	}
	if aggSMva, ok := meterMap["agg_s_mva"]; ok && aggSMva != nil {
		meter.AggregateApparentPower = *aggSMva
		meter.ApparentPowerVA = float32(*aggSMva) / 1000.0 // Convert milli VA to VA
	}
	if aggPPhAMw, ok := meterMap["agg_p_ph_a_mw"]; ok && aggPPhAMw != nil {
		meter.ActivePowerPhaseA = *aggPPhAMw
	}
	if aggPPhBMw, ok := meterMap["agg_p_ph_b_mw"]; ok && aggPPhBMw != nil {
		meter.ActivePowerPhaseB = *aggPPhBMw
	}
	if aggPPhCMw, ok := meterMap["agg_p_ph_c_mw"]; ok && aggPPhCMw != nil {
		meter.ActivePowerPhaseC = *aggPPhCMw
	}
	if aggSPhAMva, ok := meterMap["agg_s_ph_a_mva"]; ok && aggSPhAMva != nil {
		meter.ApparentPowerPhaseA = *aggSPhAMva
	}
	if aggSPhBMva, ok := meterMap["agg_s_ph_b_mva"]; ok && aggSPhBMva != nil {
		meter.ApparentPowerPhaseB = *aggSPhBMva
	}
	if aggSPhCMva, ok := meterMap["agg_s_ph_c_mva"]; ok && aggSPhCMva != nil {
		meter.ApparentPowerPhaseC = *aggSPhCMva
	}

	return meter
}

// IsSystemOnline returns true if the system is properly connected and communicating
func (ld *LiveDataStatus) IsSystemOnline() bool {
	if ld.Connection == nil {
		return false
	}
	return ld.Connection.IsConnected
}

// GetNetPowerFlow returns the net power flow (positive = importing, negative = exporting)
func (ld *LiveDataStatus) GetNetPowerFlow() float32 {
	if ld.Meters == nil || ld.Meters.Grid == nil {
		return 0
	}
	return ld.Meters.Grid.ActivePowerWatts
}

// GetStorageSOC returns the storage state of charge percentage
func (ld *LiveDataStatus) GetStorageSOC() int {
	if ld.Meters == nil {
		return 0
	}
	return ld.Meters.StateOfCharge
}
