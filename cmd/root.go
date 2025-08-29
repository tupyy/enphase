package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	verbose   bool
	rawOutput bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "enphase",
	Short: "A CLI tool for interacting with Enphase IQ Gateway Local API",
	Long: `enphase is a command-line interface for interacting with your local 
Enphase IQ Gateway device. It provides easy access to solar production data,
energy consumption information, system status, and more.

Examples:
  enphase auth login                    # Authenticate with Enphase
  enphase auth token                    # Get current auth token
  enphase production                    # Get production data
  enphase system info                   # Get system information`,
	Version: "1.0.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.enphase.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&rawOutput, "raw", true, "Output raw JSON response")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// Add all commands
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(tokenCmd)
	authCmd.AddCommand(statusCmd)

	// Login command flags
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Enphase account username (email)")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Enphase account password (will prompt if not provided)")
	loginCmd.Flags().StringVarP(&envoySerial, "envoy-serial", "s", "122312002019", "IQ Gateway serial number")
	loginCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	loginCmd.Flags().BoolVar(&saveToken, "save", true, "Save token to configuration file")

	// Mark required flags
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("envoy-serial")

	// Status command flags
	statusCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")

	rootCmd.AddCommand(ivpCmd)
	ivpCmd.AddCommand(metersCmd)
	ivpCmd.AddCommand(readingsCmd)
	ivpCmd.AddCommand(consumptionCmd)
	ivpCmd.AddCommand(gridCmd)
	ivpCmd.AddCommand(livedataCmd)
	ivpCmd.AddCommand(ensembleCmd)
	ivpCmd.AddCommand(pdmCmd)

	// Add gateway-ip flag to all ivp subcommands
	metersCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	readingsCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	consumptionCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	gridCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	livedataCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	ensembleCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	pdmCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")

	rootCmd.AddCommand(productionMainCmd)
	productionMainCmd.AddCommand(summaryCmd)
	productionMainCmd.AddCommand(invertersCmd)

	// Add gateway-ip flag to production main command and all subcommands
	productionMainCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	summaryCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
	invertersCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")

	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")

	rootCmd.AddCommand(summaryMainCmd)
	summaryMainCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")

	rootCmd.AddCommand(metricsCmd)
	metricsCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".enphase" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".enphase")
	}

	// Environment variables
	viper.SetEnvPrefix("ENPHASE")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
