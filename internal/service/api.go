package service

import (
	"encoding/json"
	"fmt"

	v1 "github.com/tupyy/enphase/api/v1"
	"github.com/tupyy/enphase/internal/client"
	"github.com/tupyy/enphase/internal/entity"
)

// ApiService wraps the client and provides typed responses for API endpoints
type ApiService struct {
	client *client.Client
}

// NewApiService creates a new API service
func NewApiService(gatewayIP string) *ApiService {
	return &ApiService{
		client: client.NewClient(gatewayIP),
	}
}

// GetProduction returns typed production data
func (s *ApiService) GetProduction() (*entity.ProductionData, error) {
	data, err := s.client.GetProduction()
	if err != nil {
		return nil, err
	}

	var production v1.ProductionData
	if err := json.Unmarshal(data, &production); err != nil {
		return nil, fmt.Errorf("failed to unmarshal production data: %w", err)
	}

	return entity.NewProductionData(&production, data), nil
}

// GetProductionInverters returns typed inverter production data
func (s *ApiService) GetProductionInverters() (*entity.InverterProduction, error) {
	data, err := s.client.GetProductionInverters()
	if err != nil {
		return nil, err
	}

	var inverters v1.InverterProduction
	if err := json.Unmarshal(data, &inverters); err != nil {
		return nil, fmt.Errorf("failed to unmarshal inverter production data: %w", err)
	}

	return entity.NewInverterProduction(&inverters, data), nil
}
