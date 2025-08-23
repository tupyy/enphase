package cmd

import (
	"encoding/json"
	"fmt"
	"os"

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

Always outputs JSON format.

Examples:
  enphase-cli summary                       # Get combined summary as JSON`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSummaryCommand()
	},
}

func init() {
	rootCmd.AddCommand(summaryMainCmd)

	// Add gateway-ip flag
	summaryMainCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
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