package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tupyy/enphase/internal/service"
)


// summaryMainCmd represents the summary command
var summaryMainCmd = &cobra.Command{
	Use:   "summary",
	Short: "Get combined consumption and production summary",
	Long: `Get a combined summary showing both current power consumption and production data.

This command fetches data from both:
- IVP consumption endpoint (cumulative current watts)
- Production inverters endpoint (sum of all inverter last report watts)
- Production data endpoint (energy today)

Outputs data in table format with the same metrics as the metrics command.

Examples:
  enphase summary                       # Get combined summary as table`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSummaryCommand()
	},
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

	// Fetch production data
	production, err := apiService.GetProduction()
	if err != nil {
		return fmt.Errorf("failed to get production data: %w", err)
	}

	// Fetch production inverters data
	inverters, err := apiService.GetProductionInverters()
	if err != nil {
		return fmt.Errorf("failed to get production inverters data: %w", err)
	}

	// Fetch PDM energy data for consumption today
	pdm, err := ivpService.GetPDM()
	if err != nil {
		return fmt.Errorf("failed to get PDM energy data: %w", err)
	}

	// Calculate total production watts from inverters (more accurate than production CT)
	var totalProductionWatts float64
	for _, inv := range inverters.Inverters {
		totalProductionWatts += float64(inv.LastReportWatts)
	}

	// Get net consumption from the grid CT
	var netConsumptionWatts float64
	for _, item := range consumption.Items {
		if item.Cumulative != nil && item.ReportType == "net-consumption" {
			netConsumptionWatts = float64(item.Cumulative.InstantaneousActivePower)
		}
	}

	// Derive total consumption: production + net (net is negative when exporting)
	totalConsumptionWatts := totalProductionWatts + netConsumptionWatts

	// Energy today from production
	energyToday := float64(production.EnergyToday)

	// Get consumption energy today from PDM
	var consumptionEnergyToday float64
	if pdm.Consumption != nil && pdm.Consumption.EIM != nil {
		consumptionEnergyToday = float64(pdm.Consumption.EIM.EnergyToday)
	}

	// Output as table
	return outputSummaryTable(netConsumptionWatts, totalConsumptionWatts, totalProductionWatts, energyToday, consumptionEnergyToday)
}

func outputSummaryTable(netConsumptionWatts, totalConsumptionWatts, productionWatts, wattHoursToday, consumptionWattHoursToday float64) error {
	// Create a new tab writer for table formatting
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	
	// Print table header
	fmt.Fprintln(w, "METRIC\tVALUE\tUNIT")
	fmt.Fprintln(w, "------\t-----\t----")
	
	// Print metrics (same as metrics command)
	fmt.Fprintf(w, "Net Consumption\t%.1f\tW\n", netConsumptionWatts)
	fmt.Fprintf(w, "Total Consumption\t%.1f\tW\n", totalConsumptionWatts)
	fmt.Fprintf(w, "Production (Current)\t%.1f\tW\n", productionWatts)
	fmt.Fprintf(w, "Production (Today)\t%.0f\tWh\n", wattHoursToday)
	fmt.Fprintf(w, "Consumption (Today)\t%.0f\tWh\n", consumptionWattHoursToday)
	
	// Flush the table
	return w.Flush()
}

