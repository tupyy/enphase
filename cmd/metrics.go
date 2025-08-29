package cmd

import (
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tupyy/enphase/internal/service"
)

// metricsCmd represents the metrics command
var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Get system metrics in Prometheus format",
	Long: `Get combined consumption and production metrics in Prometheus format.

This command fetches data from both:
- IVP consumption endpoint (cumulative current watts)
- Production inverters endpoint (sum of all inverter last report watts)

Outputs metrics in Prometheus format suitable for monitoring and alerting.

Examples:
  enphase metrics                       # Get all metrics in Prometheus format
  enphase metrics --gateway-ip 192.168.1.100  # Specify gateway IP`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMetricsCommand()
	},
}

func runMetricsCommand() error {
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

	// Get individual inverter production data
	invertersProduction := make(map[string]float64)
	for _, inv := range inverters.Inverters {
		invertersProduction[inv.SerialNumber] = float64(inv.LastReportWatts)
	}

	return outputPrometheusMetrics(currentConsumptionWatts, totalProductionWatts, invertersProduction)
}

func outputPrometheusMetrics(consumptionWatts, productionWatts float64, invertersProduction map[string]float64) error {
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

	invertersProductionGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "enphase_inverters_production_watts",
		Help: "Power production per inverters in watts",
	}, []string{"sn"})

	// Register metrics
	registry.MustRegister(consumptionGauge)
	registry.MustRegister(productionGauge)
	registry.MustRegister(realConsumptionGauge)
	registry.MustRegister(invertersProductionGauge)

	// Set values
	consumptionGauge.Set(consumptionWatts)
	productionGauge.Set(productionWatts)
	realConsumptionGauge.Set(consumptionWatts + productionWatts)

	for sn, value := range invertersProduction {
		invertersProductionGauge.With(prometheus.Labels{"sn": sn}).Set(value)
	}

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
