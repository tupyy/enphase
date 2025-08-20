package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tupyy/enphase/internal/entity"
	"github.com/tupyy/enphase/internal/presenter"
	"github.com/tupyy/enphase/internal/service"
)

// ivpCmd represents the ivp command
var ivpCmd = &cobra.Command{
	Use:   "ivp",
	Short: "Access IVP (Installer View Portal) endpoints",
	Long: `Access various IVP endpoints for meter data, live status, and device information.

IVP endpoints provide detailed information about:
- Meter status and readings
- Live data and system status  
- Energy consumption reports
- Grid connection data
- Device provisioning status

All IVP commands require authentication. Use 'enphase-cli auth login' first.

Examples:
  enphase-cli ivp meters                    # Get meter details
  enphase-cli ivp readings                  # Get meter readings
  enphase-cli ivp consumption               # Get consumption data
  enphase-cli ivp livedata                  # Get live system status`,
}

// metersCmd represents the ivp meters command
var metersCmd = &cobra.Command{
	Use:   "meters",
	Short: "Get meter details and status",
	Long: `Get meter status, type of meter, and number of phase measurements
for all installed current transformers (CTs).

This endpoint returns information about production, consumption, and storage meters.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMetersCommand()
	},
}

// readingsCmd represents the ivp readings command
var readingsCmd = &cobra.Command{
	Use:   "readings",
	Short: "Get detailed meter readings",
	Long: `Get measurements from production CT, storage CT, and consumption CT.
Data availability depends on installed CTs. Data updates every 5 minutes.

Returns detailed power measurements including:
- Active, reactive, and apparent power
- Voltage, current, and frequency
- Energy delivered and received
- Per-phase channel data`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runReadingsCommand()
	},
}

// consumptionCmd represents the ivp consumption command
var consumptionCmd = &cobra.Command{
	Use:   "consumption",
	Short: "Get power consumption reports",
	Long: `Get power consumption information of the loads including cumulative
and per-phase measurements. Data updates every 5 minutes.

Returns consumption data with:
- Net consumption vs total consumption
- Cumulative energy values
- Per-phase breakdown
- Power factor and frequency data`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runConsumptionCommand()
	},
}

// gridCmd represents the ivp grid command
var gridCmd = &cobra.Command{
	Use:   "grid",
	Short: "Get grid connection readings",
	Long: `Get voltage, current, frequency, active and reactive power
at the point of grid connection.

Negative values indicate export (feed-in), positive values indicate 
import (grid-supply). Data includes per-phase measurements.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGridCommand()
	},
}

// livedataCmd represents the ivp livedata command
var livedataCmd = &cobra.Command{
	Use:   "livedata",
	Short: "Get live system status and data",
	Long: `Get meter live data with connection status, meter readings,
task information, and system counters.

Includes real-time information about:
- MQTT connection status
- Storage state of charge
- PV, storage, grid, and load measurements
- System tasks and counters`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLivedataCommand()
	},
}

// ensembleCmd represents the ivp ensemble command
var ensembleCmd = &cobra.Command{
	Use:   "ensemble",
	Short: "Get ensemble device list",
	Long: `Get the commissioning status of provisioned devices including
storage devices, microinverters, and other connected equipment.

Returns device information including:
- Serial numbers and device types
- Connection status
- Communication interfaces
- Device capacity and phase information`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEnsembleCommand()
	},
}

// pdmCmd represents the ivp pdm command
var pdmCmd = &cobra.Command{
	Use:   "pdm",
	Short: "Get PDM energy data",
	Long: `Get energy and active power values for microinverters, revenue grade meters,
production and consumption meters for today, last seven days, and lifetime.

Works even when production meter is not installed. Includes data from:
- PCU (microinverters)
- RGM (revenue grade meters)
- EIM (Envoy internal meters)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPdmCommand()
	},
}

func init() {
	rootCmd.AddCommand(ivpCmd)

	// Add subcommands
	ivpCmd.AddCommand(metersCmd)
	ivpCmd.AddCommand(readingsCmd)
	ivpCmd.AddCommand(consumptionCmd)
	ivpCmd.AddCommand(gridCmd)
	ivpCmd.AddCommand(livedataCmd)
	ivpCmd.AddCommand(ensembleCmd)
	ivpCmd.AddCommand(pdmCmd)

	// Add gateway-ip flag to all subcommands
	metersCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	readingsCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	consumptionCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	gridCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	livedataCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	ensembleCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	pdmCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")

	// Add raw output flag to all subcommands
	metersCmd.Flags().BoolVar(&rawOutput, "raw", false, "Output raw JSON response")
	readingsCmd.Flags().BoolVar(&rawOutput, "raw", false, "Output raw JSON response")
	consumptionCmd.Flags().BoolVar(&rawOutput, "raw", false, "Output raw JSON response")
	gridCmd.Flags().BoolVar(&rawOutput, "raw", false, "Output raw JSON response")
	livedataCmd.Flags().BoolVar(&rawOutput, "raw", false, "Output raw JSON response")
	ensembleCmd.Flags().BoolVar(&rawOutput, "raw", false, "Output raw JSON response")
	pdmCmd.Flags().BoolVar(&rawOutput, "raw", false, "Output raw JSON response")
}

// Helper function to get gateway IP
func getGatewayIP() string {
	if gatewayIP == "envoy.local" {
		configIP := viper.GetString("gateway.ip")
		if configIP != "" {
			return configIP
		}
	}
	return gatewayIP
}

// Helper function to present data using appropriate presenter
func presentData(data interface{}, useRaw bool) error {
	switch v := data.(type) {
	case *entity.MeterDetails:
		p := presenter.NewMeterPresenter()
		if useRaw {
			return p.Raw(v)
		}
		return p.Summary(v)
	case *entity.MeterReadings:
		p := presenter.NewMeterReadingsPresenter()
		if useRaw {
			return p.Raw(v)
		}
		return p.Summary(v)
	case *entity.ConsumptionData:
		p := presenter.NewConsumptionPresenter()
		if useRaw {
			return p.Raw(v)
		}
		return p.Summary(v)
	case *entity.GridReading:
		p := presenter.NewGridPresenter()
		if useRaw {
			return p.Raw(v)
		}
		return p.Summary(v)
	case *entity.LiveDataStatus:
		p := presenter.NewLiveDataPresenter()
		if useRaw {
			return p.Raw(v)
		}
		return p.Summary(v)
	case *entity.DeviceList:
		p := presenter.NewDevicePresenter()
		if useRaw {
			return p.Raw(v)
		}
		return p.Summary(v)
	case *entity.EnergyData:
		p := presenter.NewEnergyPresenter()
		if useRaw {
			return p.Raw(v)
		}
		return p.Summary(v)
	case *entity.ProductionData:
		p := presenter.NewProductionPresenter()
		if useRaw {
			return p.Raw(v)
		}
		return p.Summary(v)
	case *entity.InverterProduction:
		p := presenter.NewInverterPresenter()
		if useRaw {
			return p.Raw(v)
		}
		return p.Summary(v)
	default:
		return fmt.Errorf("unsupported data type: %T", data)
	}
}

func runMetersCommand() error {
	ip := getGatewayIP()
	if verbose {
		fmt.Printf("Connecting to gateway: %s\n", ip)
		fmt.Printf("Endpoint: /ivp/meters\n")
	}

	ivpService := service.NewIvpService(ip)
	meters, err := ivpService.GetMeters()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Response received (meters data)\n")
		fmt.Println("---")
	}

	return presentData(meters, rawOutput)
}

func runReadingsCommand() error {
	ip := getGatewayIP()
	if verbose {
		fmt.Printf("Connecting to gateway: %s\n", ip)
		fmt.Printf("Endpoint: /ivp/meters/readings\n")
	}

	ivpService := service.NewIvpService(ip)
	readings, err := ivpService.GetReadings()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Response received (readings data)\n")
		fmt.Println("---")
	}

	return presentData(readings, rawOutput)
}

func runConsumptionCommand() error {
	ip := getGatewayIP()
	if verbose {
		fmt.Printf("Connecting to gateway: %s\n", ip)
		fmt.Printf("Endpoint: /ivp/meters/reports/consumption\n")
	}

	ivpService := service.NewIvpService(ip)
	consumption, err := ivpService.GetConsumption()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Response received (consumption data)\n")
		fmt.Println("---")
	}

	return presentData(consumption, rawOutput)
}

func runGridCommand() error {
	ip := getGatewayIP()
	if verbose {
		fmt.Printf("Connecting to gateway: %s\n", ip)
		fmt.Printf("Endpoint: /ivp/meters/gridReading\n")
	}

	ivpService := service.NewIvpService(ip)
	grid, err := ivpService.GetGrid()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Response received (grid data)\n")
		fmt.Println("---")
	}

	return presentData(grid, rawOutput)
}

func runLivedataCommand() error {
	ip := getGatewayIP()
	if verbose {
		fmt.Printf("Connecting to gateway: %s\n", ip)
		fmt.Printf("Endpoint: /ivp/livedata/status\n")
	}

	ivpService := service.NewIvpService(ip)
	livedata, err := ivpService.GetLiveData()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Response received (livedata)\n")
		fmt.Println("---")
	}

	return presentData(livedata, rawOutput)
}

func runEnsembleCommand() error {
	ip := getGatewayIP()
	if verbose {
		fmt.Printf("Connecting to gateway: %s\n", ip)
		fmt.Printf("Endpoint: /ivp/ensemble/device_list\n")
	}

	ivpService := service.NewIvpService(ip)
	ensemble, err := ivpService.GetEnsemble()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Response received (ensemble data)\n")
		fmt.Println("---")
	}

	return presentData(ensemble, rawOutput)
}

func runPdmCommand() error {
	ip := getGatewayIP()
	if verbose {
		fmt.Printf("Connecting to gateway: %s\n", ip)
		fmt.Printf("Endpoint: /ivp/pdm/energy\n")
	}

	ivpService := service.NewIvpService(ip)
	pdm, err := ivpService.GetPDM()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Response received (PDM energy data)\n")
		fmt.Println("---")
	}

	return presentData(pdm, rawOutput)
}
