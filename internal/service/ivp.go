package service

import (
	"encoding/json"
	"fmt"

	v1 "github.com/tupyy/enphase/api/v1"
	"github.com/tupyy/enphase/internal/client"
	"github.com/tupyy/enphase/internal/entity"
)

// IvpService wraps the client and provides typed responses for IVP endpoints
type IvpService struct {
	client *client.Client
}

// NewIvpService creates a new IVP service
func NewIvpService(gatewayIP string) *IvpService {
	return &IvpService{
		client: client.NewClient(gatewayIP),
	}
}

// GetMeters returns typed meter details
func (s *IvpService) GetMeters() (*entity.MeterDetails, error) {
	data, err := s.client.GetMeters()
	if err != nil {
		return nil, err
	}

	var meters v1.MeterDetails
	if err := json.Unmarshal(data, &meters); err != nil {
		return nil, fmt.Errorf("failed to unmarshal meter details: %w", err)
	}

	return entity.NewMeterDetails(&meters, data), nil
}

// GetReadings returns typed meter readings
func (s *IvpService) GetReadings() (*entity.MeterReadings, error) {
	data, err := s.client.GetMeterReadings()
	if err != nil {
		return nil, err
	}

	var readings v1.MeterReadings
	if err := json.Unmarshal(data, &readings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal meter readings: %w", err)
	}

	return entity.NewMeterReadings(&readings, data), nil
}

// GetConsumption returns typed consumption data
func (s *IvpService) GetConsumption() (*entity.ConsumptionData, error) {
	data, err := s.client.GetMeterConsumption()
	if err != nil {
		return nil, err
	}

	// Try to unmarshal as array first (which is what the API returns)
	var consumptionArray []v1.ConsumptionData
	if err := json.Unmarshal(data, &consumptionArray); err != nil {
		// If that fails, try as single object (fallback)
		var consumption v1.ConsumptionData
		if err := json.Unmarshal(data, &consumption); err != nil {
			return nil, fmt.Errorf("failed to unmarshal consumption data: %w", err)
		}
		// Convert single object to array for consistent processing
		consumptionArray = []v1.ConsumptionData{consumption}
	}

	// Return all elements in the array
	return entity.NewConsumptionData(consumptionArray, data), nil
}

// GetGrid returns typed grid reading data
func (s *IvpService) GetGrid() (*entity.GridReading, error) {
	data, err := s.client.GetMeterGridReading()
	if err != nil {
		return nil, err
	}

	var grid v1.GridReading
	if err := json.Unmarshal(data, &grid); err != nil {
		return nil, fmt.Errorf("failed to unmarshal grid reading: %w", err)
	}

	return entity.NewGridReading(&grid, data), nil
}

// GetLiveData returns typed live data status
func (s *IvpService) GetLiveData() (*entity.LiveDataStatus, error) {
	data, err := s.client.GetLiveDataStatus()
	if err != nil {
		return nil, err
	}

	var livedata v1.LiveDataStatus
	if err := json.Unmarshal(data, &livedata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal live data status: %w", err)
	}

	return entity.NewLiveDataStatus(&livedata, data), nil
}

// GetEnsemble returns typed device list
func (s *IvpService) GetEnsemble() (*entity.DeviceList, error) {
	data, err := s.client.GetEnsembleDeviceList()
	if err != nil {
		return nil, err
	}

	var ensemble v1.DeviceList
	if err := json.Unmarshal(data, &ensemble); err != nil {
		return nil, fmt.Errorf("failed to unmarshal device list: %w", err)
	}

	return entity.NewDeviceList(&ensemble, data), nil
}

// GetPDM returns typed PDM energy data
func (s *IvpService) GetPDM() (*entity.EnergyData, error) {
	data, err := s.client.GetPDMEnergy()
	if err != nil {
		return nil, err
	}

	var pdm v1.EnergyData
	if err := json.Unmarshal(data, &pdm); err != nil {
		return nil, fmt.Errorf("failed to unmarshal PDM energy data: %w", err)
	}

	return entity.NewEnergyData(&pdm, data), nil
}
