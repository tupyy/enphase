package client

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

// Client represents a unified HTTP client for the Enphase IQ Gateway
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClient creates a new unified client
func NewClient(gatewayIP string) *Client {
	return &Client{
		BaseURL: fmt.Sprintf("https://%s", gatewayIP),
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		Token: viper.GetString("auth.token"),
	}
}

// makeRequest makes an authenticated request to the gateway
func (c *Client) makeRequest(endpoint string) ([]byte, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("no authentication token found. Please run 'enphase-cli auth login' first")
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gateway: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("authentication failed. Token may be expired. Please run 'enphase-cli auth login'")
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("endpoint not available. Feature may not be installed or enabled")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gateway returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// IVP Endpoints

// GetMeters gets meter details
func (c *Client) GetMeters() ([]byte, error) {
	return c.makeRequest("/ivp/meters")
}

// GetMeterReadings gets meter readings
func (c *Client) GetMeterReadings() ([]byte, error) {
	return c.makeRequest("/ivp/meters/readings")
}

// GetMeterConsumption gets consumption reports
func (c *Client) GetMeterConsumption() ([]byte, error) {
	return c.makeRequest("/ivp/meters/reports/consumption")
}

// GetMeterGridReading gets grid readings
func (c *Client) GetMeterGridReading() ([]byte, error) {
	return c.makeRequest("/ivp/meters/gridReading")
}

// GetLiveDataStatus gets live data status
func (c *Client) GetLiveDataStatus() ([]byte, error) {
	return c.makeRequest("/ivp/livedata/status")
}

// GetEnsembleDeviceList gets ensemble device list
func (c *Client) GetEnsembleDeviceList() ([]byte, error) {
	return c.makeRequest("/ivp/ensemble/device_list")
}

// GetPDMEnergy gets PDM energy data
func (c *Client) GetPDMEnergy() ([]byte, error) {
	return c.makeRequest("/ivp/pdm/energy")
}

// API Endpoints

// GetProduction gets production meter data
func (c *Client) GetProduction() ([]byte, error) {
	return c.makeRequest("/api/v1/production")
}

// GetProductionInverters gets inverter production data
func (c *Client) GetProductionInverters() ([]byte, error) {
	return c.makeRequest("/api/v1/production/inverters")
}
