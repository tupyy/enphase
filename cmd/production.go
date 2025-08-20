package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tupyy/enphase/internal/service"
)

// productionMainCmd represents the production command
var productionMainCmd = &cobra.Command{
	Use:   "production",
	Short: "Access production data and inverter information",
	Long: `Access production data and inverter information from your IQ Gateway.

Production endpoints provide access to:
- Production meter data and energy summaries
- Individual microinverter production data
- System-wide energy statistics

All production commands require authentication. Use 'enphase-cli auth login' first.

Examples:
  enphase-cli production                    # Get production summary (default)
  enphase-cli production inverters          # Get inverter details`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default action: get production summary
		return runProductionSummary()
	},
}

// summaryCmd represents the production summary command
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Get production meter summary data",
	Long: `Get production energy and active power values for today, last seven days,
and lifetime in watt-hours. 

This endpoint only works when a production meter is installed and enabled at the site.

Returns:
- Today's energy production (Wh)
- Last 7 days energy production (Wh) 
- Lifetime energy production (Wh)
- Current power production (W)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runProductionSummary()
	},
}

// invertersCmd represents the production inverters command
var invertersCmd = &cobra.Command{
	Use:   "inverters",
	Short: "Get inverter production data",
	Long: `Get maximum and last reported active power production information
for all available microinverters. Data updates every 5 minutes.

Returns for each inverter:
- Serial number
- Last report date/time
- Device type
- Last reported power (Watts)
- Maximum ever reported power (Watts)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runProductionInverters()
	},
}

func init() {
	rootCmd.AddCommand(productionMainCmd)

	// Add subcommands
	productionMainCmd.AddCommand(summaryCmd)
	productionMainCmd.AddCommand(invertersCmd)

	// Add gateway-ip flag to main command and all subcommands
	productionMainCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	summaryCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	invertersCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")

	// Add raw output flag to main command and all subcommands
	productionMainCmd.Flags().BoolVar(&rawOutput, "raw", false, "Output raw JSON response")
	summaryCmd.Flags().BoolVar(&rawOutput, "raw", false, "Output raw JSON response")
	invertersCmd.Flags().BoolVar(&rawOutput, "raw", false, "Output raw JSON response")
}

func runProductionSummary() error {
	// Get gateway IP from config if not specified
	if gatewayIP == "envoy.local" {
		configIP := viper.GetString("gateway.ip")
		if configIP != "" {
			gatewayIP = configIP
		}
	}

	if verbose {
		fmt.Printf("Connecting to gateway: %s\n", gatewayIP)
		fmt.Printf("Endpoint: /api/v1/production\n")
	}

	// Create API service
	apiService := service.NewApiService(gatewayIP)

	production, err := apiService.GetProduction()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Response received (production data)\n")
		fmt.Println("---")
	}

	// Present using appropriate presenter
	return presentData(production, rawOutput)
}

func runProductionInverters() error {
	// Get gateway IP from config if not specified
	if gatewayIP == "envoy.local" {
		configIP := viper.GetString("gateway.ip")
		if configIP != "" {
			gatewayIP = configIP
		}
	}

	if verbose {
		fmt.Printf("Connecting to gateway: %s\n", gatewayIP)
		fmt.Printf("Endpoint: /api/v1/production/inverters\n")
	}

	// Create API service
	apiService := service.NewApiService(gatewayIP)

	inverters, err := apiService.GetProductionInverters()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Response received (inverter data)\n")
		fmt.Println("---")
	}

	// Present using appropriate presenter
	return presentData(inverters, rawOutput)
}
