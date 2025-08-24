package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tupyy/enphase/internal/service"
)

var (
	metricsOutput bool
)

// summaryMainCmd represents the summary command
var summaryMainCmd = &cobra.Command{
	Use:   "summary",
	Short: "Get combined consumption and production summary",
	Long: `Get a combined summary showing both current power consumption and production data.

This command fetches data from both:
- IVP consumption endpoint (cumulative current watts)
- Production inverters endpoint (sum of all inverter last report watts)

Outputs JSON format by default, or Prometheus metrics format with --metrics flag.

Examples:
  enphase-cli summary                       # Get combined summary as JSON
  enphase-cli summary --metrics             # Get combined summary as Prometheus metrics`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSummaryCommand()
	},
}

func init() {
	rootCmd.AddCommand(summaryMainCmd)

	// Add gateway-ip flag
	summaryMainCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	// Add metrics flag
	summaryMainCmd.Flags().BoolVar(&metricsOutput, "metrics", false, "Output in Prometheus metrics format")
}

func runSummaryCommand() error {
	// Get gateway IP from config if not specified
	if gatewayIP == "envoy.local" {
		configIP := viper.GetString("gateway.ip")
		if configIP != "" {
			gatewayIP = configIP
		}
	}

	// Create services
	ivpService := service.NewIvpService(gatewayIP)
	apiService := service.NewApiService(gatewayIP)

	// Fetch consumption data
	consumption, err := ivpService.GetConsumption()
	if err != nil {
		return fmt.Errorf("failed to get consumption data: %w", err)
	}

	// Fetch production inverters data
	inverters, err := apiService.GetProductionInverters()
	if err != nil {
		return fmt.Errorf("failed to get production inverters data: %w", err)
	}

	// Calculate total production watts
	totalProductionWatts := float64(inverters.TotalCurrentPower())

	// Get current consumption watts from cumulative data
	var currentConsumptionWatts float64
	if consumption.Cumulative != nil {
		currentConsumptionWatts = float64(consumption.Cumulative.InstantaneousActivePower)
	}

	if metricsOutput {
		return outputPrometheusMetrics(currentConsumptionWatts, totalProductionWatts)
	}

	// Format to 1 decimal place
	formattedSummary := struct {
		ConsumptionWatts string `json:"consumption_watts"`
		ProductionWatts  string `json:"production_watts"`
	}{
		ConsumptionWatts: fmt.Sprintf("%.1f", currentConsumptionWatts),
		ProductionWatts:  fmt.Sprintf("%.1f", totalProductionWatts),
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(formattedSummary)
}

func outputPrometheusMetrics(consumptionWatts, productionWatts float64) error {
	// Create a new registry
	registry := prometheus.NewRegistry()

	// Create metrics
	consumptionGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "enphase_consumption_watts",
		Help: "Current power consumption in watts (net consumption from grid)",
	})

	productionGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "enphase_production_watts",
		Help: "Current power production in watts",
	})

	realConsumptionGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "enphase_real_consumption_watts",
		Help: "Real household power consumption in watts (consumption + production)",
	})

	// Register metrics
	registry.MustRegister(consumptionGauge)
	registry.MustRegister(productionGauge)
	registry.MustRegister(realConsumptionGauge)

	// Set values
	consumptionGauge.Set(consumptionWatts)
	productionGauge.Set(productionWatts)
	realConsumptionGauge.Set(consumptionWatts + productionWatts)

	// Gather metrics
	metricFamilies, err := registry.Gather()
	if err != nil {
		return fmt.Errorf("failed to gather metrics: %w", err)
	}

	// Write metrics to stdout in Prometheus format
	for _, mf := range metricFamilies {
		if _, err := expfmt.MetricFamilyToText(os.Stdout, mf); err != nil {
			return fmt.Errorf("failed to write metric family: %w", err)
		}
	}

	return nil
}