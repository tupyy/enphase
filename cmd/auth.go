package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var (
	username    string
	password    string
	envoySerial string
	gatewayIP   string
	saveToken   bool
	rawOutput   bool
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication with Enphase IQ Gateway",
	Long: `Manage authentication tokens for accessing the Enphase IQ Gateway Local API.

This command helps you obtain and manage Bearer tokens required for accessing
protected endpoints on your IQ Gateway device.`,
}

// loginCmd represents the auth login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate and obtain a Bearer token",
	Long: `Authenticate with Enphase and obtain a Bearer token for API access.

This command will:
1. Authenticate with Enphase Cloud using your credentials
2. Generate a Bearer token for your specific IQ Gateway
3. Optionally save the token to your configuration file

Examples:
  enphase-cli auth login --username user@example.com --envoy-serial 123456789012
  enphase-cli auth login -u user@example.com -s 123456789012 --save`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performLogin()
	},
}

// tokenCmd represents the auth token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Display the current Bearer token",
	Long: `Display the current Bearer token stored in configuration.

If no token is stored, you will need to run 'enphase-cli auth login' first.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showToken()
	},
}

// statusCmd represents the auth status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",
	Long: `Check if the current token is valid by testing it against the IQ Gateway.

This command will attempt to access a protected endpoint to verify token validity.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return checkAuthStatus()
	},
}

func init() {
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
}

func performLogin() error {
	// Get password if not provided
	if password == "" {
		fmt.Print("Enter password: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		password = string(bytePassword)
		fmt.Println() // New line after password input
	}

	if verbose {
		fmt.Printf("Authenticating user: %s\n", username)
		fmt.Printf("Envoy serial: %s\n", envoySerial)
	}

	// Step 1: Get session ID
	sessionID, err := getSessionID()
	if err != nil {
		return fmt.Errorf("failed to get session ID: %w", err)
	}

	if verbose {
		fmt.Printf("Session ID obtained: %s\n", sessionID)
	}

	// Step 2: Get Bearer token
	token, err := getBearerToken(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get Bearer token: %w", err)
	}

	fmt.Printf("✅ Authentication successful!\n")
	fmt.Printf("Bearer Token: %s\n", token)

	// Save token if requested
	if saveToken {
		viper.Set("auth.token", token)
		viper.Set("auth.username", username)
		viper.Set("auth.envoy_serial", envoySerial)
		viper.Set("gateway.ip", gatewayIP)

		if err := viper.WriteConfig(); err != nil {
			// If config file doesn't exist, create it
			if err := viper.SafeWriteConfig(); err != nil {
				fmt.Printf("⚠️  Warning: Could not save token to config: %v\n", err)
			} else {
				fmt.Printf("✅ Token saved to configuration file\n")
			}
		} else {
			fmt.Printf("✅ Token saved to configuration file\n")
		}
	}

	return nil
}

func getSessionID() (string, error) {
	// Prepare form data
	data := url.Values{}
	data.Set("user[email]", username)
	data.Set("user[password]", password)

	// Make POST request
	resp, err := http.PostForm("https://enlighten.enphaseenergy.com/login/login.json", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse login response: %w", err)
	}

	sessionID, ok := result["session_id"].(string)
	if !ok {
		return "", fmt.Errorf("session_id not found in response")
	}

	return sessionID, nil
}

func getBearerToken(sessionID string) (string, error) {
	// Prepare JSON payload
	payload := map[string]string{
		"session_id": sessionID,
		"serial_num": envoySerial,
		"username":   username,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Make POST request
	req, err := http.NewRequest("POST", "https://entrez.enphaseenergy.com/tokens", strings.NewReader(string(jsonData)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// The response should be the raw token
	token := strings.TrimSpace(string(body))
	if token == "" {
		return "", fmt.Errorf("empty token received")
	}

	return token, nil
}

func showToken() error {
	token := viper.GetString("auth.token")
	if token == "" {
		fmt.Println("❌ No token found. Please run 'enphase-cli auth login' first.")
		return nil
	}

	fmt.Printf("Current Bearer Token: %s\n", token)

	// Show additional info if available
	username := viper.GetString("auth.username")
	envoySerial := viper.GetString("auth.envoy_serial")
	gatewayIP := viper.GetString("gateway.ip")

	if username != "" {
		fmt.Printf("Username: %s\n", username)
	}
	if envoySerial != "" {
		fmt.Printf("Envoy Serial: %s\n", envoySerial)
	}
	if gatewayIP != "" {
		fmt.Printf("Gateway IP: %s\n", gatewayIP)
	}

	return nil
}

func checkAuthStatus() error {
	token := viper.GetString("auth.token")
	if token == "" {
		fmt.Println("❌ No token found. Please run 'enphase-cli auth login' first.")
		return nil
	}

	// Get gateway IP from config or flag
	if gatewayIP == "envoy.local" {
		configIP := viper.GetString("gateway.ip")
		if configIP != "" {
			gatewayIP = configIP
		}
	}

	if verbose {
		fmt.Printf("Testing token against gateway: %s\n", gatewayIP)
	}

	// Test token by accessing a protected endpoint
	url := fmt.Sprintf("https://%s/api/v1/production/inverters", gatewayIP)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")

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

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Printf("✅ Token is valid and working\n")
		fmt.Printf("Successfully connected to gateway at: %s\n", gatewayIP)
	case http.StatusUnauthorized:
		fmt.Printf("❌ Token is invalid or expired\n")
		fmt.Printf("Please run 'enphase-cli auth login' to get a new token\n")
	default:
		fmt.Printf("⚠️  Unexpected response from gateway: %d\n", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 && verbose {
			fmt.Printf("Response: %s\n", string(body))
		}
	}

	return nil
}
