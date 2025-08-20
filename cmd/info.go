package cmd

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// EnvoyInfo represents the XML structure returned by /info endpoint
type EnvoyInfo struct {
	XMLName  xml.Name  `xml:"envoy_info"`
	Time     int64     `xml:"time"`
	Device   Device    `xml:"device"`
	WebToken bool      `xml:"web-tokens"`
	Packages []Package `xml:"package"`
}

type Device struct {
	SerialNumber string `xml:"sn"`
	PartNumber   string `xml:"pn"`
	Software     string `xml:"software"`
	EUAID        string `xml:"euaid"`
	SeqNum       int    `xml:"seqnum"`
	APIVersion   int    `xml:"apiver"`
	IMeter       bool   `xml:"imeter"`
}

type Package struct {
	Name    string `xml:"name,attr"`
	PN      string `xml:"pn"`
	Version string `xml:"version"`
	Build   string `xml:"build"`
}

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get system information from IQ Gateway",
	Long: `Get system information from the IQ Gateway including serial number,
model ID, firmware version, and installed packages.

This endpoint does not require authentication.

Examples:
  enphase-cli info                              # Get system info
  enphase-cli info --gateway-ip 192.168.1.100  # Use specific IP`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return getSystemInfo()
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Info command flags
	infoCmd.Flags().StringVarP(&gatewayIP, "gateway-ip", "g", "envoy.local", "IQ Gateway IP address or hostname")
}

func getSystemInfo() error {
	// Get gateway IP from config if not specified
	if gatewayIP == "envoy.local" {
		configIP := viper.GetString("gateway.ip")
		if configIP != "" {
			gatewayIP = configIP
		}
	}

	if verbose {
		fmt.Printf("Getting system info from gateway: %s\n", gatewayIP)
	}

	// Create request
	url := fmt.Sprintf("https://%s/info", gatewayIP)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Create client that ignores SSL certificate errors (self-signed certificates)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to gateway: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("gateway returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if verbose {
		fmt.Printf("Raw response:\n%s\n\n", string(body))
	}

	// Parse XML response
	var envoyInfo EnvoyInfo
	err = xml.Unmarshal(body, &envoyInfo)
	if err != nil {
		return fmt.Errorf("failed to parse XML response: %w", err)
	}

	// Display system information
	displaySystemInfo(&envoyInfo)

	return nil
}

func displaySystemInfo(info *EnvoyInfo) {
	fmt.Printf("🔌 IQ Gateway System Information\n")
	fmt.Printf("================================\n\n")

	// Device Information
	fmt.Printf("📱 Device Details:\n")
	fmt.Printf("  Serial Number:    %s\n", info.Device.SerialNumber)
	fmt.Printf("  Part Number:      %s\n", info.Device.PartNumber)
	fmt.Printf("  Software Version: %s\n", info.Device.Software)
	fmt.Printf("  EUAID:            %s\n", info.Device.EUAID)
	fmt.Printf("  API Version:      %d\n", info.Device.APIVersion)
	fmt.Printf("  Internal Meter:   %t\n", info.Device.IMeter)
	fmt.Printf("  Web Tokens:       %t\n", info.WebToken)
	fmt.Printf("\n")

	// Package Information
	fmt.Printf("📦 Installed Packages:\n")
	for _, pkg := range info.Packages {
		fmt.Printf("  %-12s v%-12s (build: %s)\n",
			pkg.Name, pkg.Version, pkg.Build)
	}
	fmt.Printf("\n")

	// Additional info
	fmt.Printf("ℹ️  Additional Info:\n")
	fmt.Printf("  Last Update:      %d (epoch timestamp)\n", info.Time)
	fmt.Printf("  Sequence Number:  %d\n", info.Device.SeqNum)
}
