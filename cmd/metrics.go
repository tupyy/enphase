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

	// Fetch meter details to identify meter types
	meters, err := ivpService.GetMeters()
	if err != nil {
		return fmt.Errorf("failed to get meter details: %w", err)
	}

	// Fetch meter readings for net-consumption data
	readings, err := ivpService.GetReadings()
	if err != nil {
		return fmt.Errorf("failed to get meter readings: %w", err)
	}

	// Fetch production inverters data
	inverters, err := apiService.GetProductionInverters()
	if err != nil {
		return fmt.Errorf("failed to get production inverters data: %w", err)
	}

	// Sum inverter production (more accurate than the miscalibrated production CT)
	var totalProductionWatts float64
	invertersProduction := make(map[string]float64)
	for _, inv := range inverters.Inverters {
		invertersProduction[inv.SerialNumber] = float64(inv.LastReportWatts)
		totalProductionWatts += float64(inv.LastReportWatts)
	}

	// Find net-consumption meter eid
	var netConsumptionEid int64
	for _, m := range meters.Meters {
		if m.IsConsumption {
			netConsumptionEid = m.ID
			break
		}
	}

	// Get net consumption from meter readings
	var netConsumptionWatts float64
	var netConsumptionActivePower float64
	for _, r := range readings.Readings {
		if r.ID == netConsumptionEid {
			netConsumptionWatts = float64(r.InstantaneousDemand)
			netConsumptionActivePower = float64(r.ActivePower)
			break
		}
	}

	// Derive total consumption: production + net (net is negative when exporting)
	totalConsumptionWatts := totalProductionWatts + netConsumptionWatts

	return outputPrometheusMetrics(netConsumptionWatts, netConsumptionActivePower, totalConsumptionWatts, totalProductionWatts, invertersProduction)
}

func outputPrometheusMetrics(netConsumptionWatts, netConsumptionActivePower, totalConsumptionWatt, productionWatts float64, invertersProduction map[string]float64) error {
	registry := prometheus.NewRegistry()

	consumptionGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "enphase_net_consumption_watts",
		Help: "Current power consumption in watts (net consumption from grid)",
	})

	consumptionActivePowerGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "enphase_net_consumption_active_power_watts",
		Help: "Active power net consumption in watts (smoothed, from grid CT)",
	})

	productionGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "enphase_production_watts",
		Help: "Current power production in watts (sum of inverters)",
	})

	totalConsumptionGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "enphase_total_consumption_watts",
		Help: "Total power consumption in watts (production + net)",
	})

	invertersProductionGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "enphase_inverters_production_watts",
		Help: "Power production per inverter in watts",
	}, []string{"sn"})

	registry.MustRegister(consumptionActivePowerGauge)
	registry.MustRegister(consumptionGauge)
	registry.MustRegister(productionGauge)
	registry.MustRegister(totalConsumptionGauge)
	registry.MustRegister(invertersProductionGauge)

	consumptionActivePowerGauge.Set(netConsumptionActivePower)
	consumptionGauge.Set(netConsumptionWatts)
	productionGauge.Set(productionWatts)
	totalConsumptionGauge.Set(totalConsumptionWatt)

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
