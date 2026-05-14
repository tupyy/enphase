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
	production, err := apiService.GetProduction()
	if err != nil {
		return fmt.Errorf("failed to get production data: %w", err)
	}

	inverters, err := apiService.GetProductionInverters()
	if err != nil {
		return fmt.Errorf("failed to get production inverters data: %w", err)
	}

	// Fetch PDM energy data for consumption today
	pdm, err := ivpService.GetPDM()
	if err != nil {
		return fmt.Errorf("failed to get PDM energy data: %w", err)
	}

	// Calibrated production CT (raw CT reads ~56.26% of actual, per regression against inverter sum)
	totalProductionWatts := float64(production.PowerNow) / 0.5626

	// Get net consumption from the grid CT
	var netConsumptionWatts float64
	var netConsumptionActivePower float64
	for _, item := range consumption.Items {
		if item.Cumulative != nil && item.ReportType == "net-consumption" {
			netConsumptionWatts = float64(item.Cumulative.InstantaneousActivePower)
			netConsumptionActivePower = float64(item.Cumulative.ActivePower)
		}
	}

	// Derive total consumption: production + net (net is negative when exporting)
	totalConsumptionWatts := totalProductionWatts + netConsumptionWatts

	// Get consumption energy today from PDM
	var consumptionEnergyToday float64
	if pdm.Consumption != nil && pdm.Consumption.EIM != nil {
		consumptionEnergyToday = float64(pdm.Consumption.EIM.EnergyToday)
	}

	// Get individual inverter production data
	invertersProduction := make(map[string]float64)
	for _, inv := range inverters.Inverters {
		invertersProduction[inv.SerialNumber] = float64(inv.LastReportWatts)
	}

	return outputPrometheusMetrics(netConsumptionWatts, netConsumptionActivePower, totalConsumptionWatts, totalProductionWatts, float64(production.EnergyToday), consumptionEnergyToday, invertersProduction)
}

func outputPrometheusMetrics(netConsumptionWatts, netConsumptionActivePower, totalConsumptionWatt, productionWatts, wattHoursToday, consumptionWattHoursToday float64, invertersProduction map[string]float64) error {
	// Create a new registry
	registry := prometheus.NewRegistry()

	// Create metrics
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
		Help: "Current power production in watts",
	})

	totalConsumptionGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "enphase_total_consumption_watts",
		Help: "Total power consumption in watts",
	})

	invertersProductionGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "enphase_inverters_production_watts",
		Help: "Power production per inverters in watts",
	}, []string{"sn"})

	wattHoursTodayGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "enphase_production_watt_hours_today",
		Help: "Total energy production today in watt-hours",
	})

	consumptionWattHoursTodayGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "enphase_consumption_watt_hours_today",
		Help: "Total energy consumption today in watt-hours",
	})

	// Register metrics
	registry.MustRegister(consumptionActivePowerGauge)
	registry.MustRegister(consumptionGauge)
	registry.MustRegister(productionGauge)
	registry.MustRegister(totalConsumptionGauge)
	registry.MustRegister(invertersProductionGauge)
	registry.MustRegister(wattHoursTodayGauge)
	registry.MustRegister(consumptionWattHoursTodayGauge)

	// Set values
	consumptionActivePowerGauge.Set(netConsumptionActivePower)
	consumptionGauge.Set(netConsumptionWatts)
	productionGauge.Set(productionWatts)
	totalConsumptionGauge.Set(totalConsumptionWatt)
	wattHoursTodayGauge.Set(wattHoursToday)
	consumptionWattHoursTodayGauge.Set(consumptionWattHoursToday)

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
