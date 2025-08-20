# Enphase IQ Gateway Local API

A comprehensive OpenAPI specification and CLI tool for the Enphase IQ Gateway Local API, providing access to solar energy production, consumption, and system data directly from your local IQ Gateway device.

**Disclaimer**: This is an unofficial OpenAPI specification created from public Enphase documentation. Always refer to official Enphase documentation for the most current information.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Authentication](#authentication)
- [Quick Start](#quick-start)
- [API Endpoints](#api-endpoints)
- [Usage Examples](#usage-examples)
- [Hardware Requirements](#hardware-requirements)
- [Token Management](#token-management)
- [Data Update Frequencies](#data-update-frequencies)
- [Contributing](#contributing)
- [License](#license)

## 🔍 Overview

This project provides a complete OpenAPI 3.0 specification for the Enphase IQ Gateway Local API, extracted from the official Enphase technical documentation. The API allows you to access real-time and historical data from your solar energy system including:

- **Energy Production**: Solar panel and microinverter data
- **Energy Consumption**: Load and consumption measurements  
- **System Information**: Device details and firmware versions
- **Grid Connection**: Voltage, current, and power flow data
- **Storage Systems**: Battery status and charge levels

## ✨ Features

- **Complete API Coverage**: All 10 documented endpoints with detailed schemas
- **Command-Line Interface**: Easy-to-use CLI tool built with Go and Cobra
- **Flexible Output Formats**: Summary (human-readable) and raw JSON output modes
- **Authentication Ready**: Bearer token authentication implementation
- **Rich Documentation**: Comprehensive descriptions and examples
- **Developer Friendly**: Ready for code generation and integration
- **Real-world Tested**: Based on actual Enphase documentation and responses
- **Domain Entities**: Type-safe data models with business logic methods
- **Presenter Layer**: Clean separation of data formatting from core logic

## 🔐 Authentication

The API uses Bearer token authentication. Tokens can be obtained through:

### Web Portal Method
1. Visit https://entrez.enphaseenergy.com
2. Login with your Enphase account
3. Select your system and IQ Gateway
4. Generate an access token

### Programmatic Method
```bash
# Get session ID
SESSION_ID=$(curl -X POST http://enlighten.enphaseenergy.com/login/login.json? \
  -F "user[email]=$USERNAME" \
  -F "user[password]=$PASSWORD" | jq -r ".session_id")

# Get token
TOKEN=$(curl -X POST http://entrez.enphaseenergy.com/tokens \
  -H "Content-Type: application/json" \
  -d "{\"session_id\": \"$SESSION_ID\", \"serial_num\": \"$ENVOY_SERIAL\", \"username\": \"$USERNAME\"}")
```

### Token Validity
- **System Owner**: 1 year
- **Installer**: 12 hours

## 🚀 Quick Start

### 1. Build the CLI Tool
```bash
git clone <repository>
cd enphase
make build
```

### 2. Find Your IQ Gateway IP
```bash
# Usually discoverable as:
ping envoy.local
# or check your router's DHCP table
```

### 3. Test Connection (No Auth Required)
```bash
./bin/enphase-cli info
```

### 4. Authenticate
```bash
./bin/enphase-cli auth login --username user@example.com --envoy-serial 123456789012
```

### 5. Get Production Data
```bash
# Get production summary
./bin/enphase-cli production

# Get inverter details
./bin/enphase-cli production inverters

# Get live system status
./bin/enphase-cli ivp livedata
```

### 6. Use the OpenAPI Spec
- **Documentation**: Load `api/v1/openapi.yaml` in Swagger UI
- **Code Generation**: Use with OpenAPI generators
- **Testing**: Import into Postman or Insomnia

## 📊 API Endpoints

| Endpoint | Method | Description | Auth Required |
|----------|--------|-------------|---------------|
| `/info` | GET | System information and firmware versions | ❌ |
| `/ivp/ensemble/device_list` | GET | Provisioned devices and status | ✅ |
| `/ivp/meters` | GET | Meter details and configuration | ✅ |
| `/ivp/meters/readings` | GET | Detailed meter measurements | ✅ |
| `/api/v1/production` | GET | Production energy summary | ✅ |
| `/ivp/pdm/energy` | GET | Comprehensive energy data | ✅ |
| `/api/v1/production/inverters` | GET | Individual inverter data | ✅ |
| `/ivp/livedata/status` | GET | Real-time system status | ✅ |
| `/ivp/meters/reports/consumption` | GET | Consumption reports | ✅ |
| `/ivp/meters/gridReading` | GET | Grid connection readings | ✅ |

## 💡 CLI Usage Examples

### System Information
```bash
# Get system info (no authentication required)
./bin/enphase-cli info

# Use specific gateway IP
./bin/enphase-cli info --gateway-ip 192.168.1.100
```

### Authentication
```bash
# Login and save token
./bin/enphase-cli auth login --username user@example.com --envoy-serial 123456789012

# Check current token
./bin/enphase-cli auth token

# Test token validity
./bin/enphase-cli auth status
```

### Production Data
```bash
# Get production summary
./bin/enphase-cli production

# Get detailed inverter data
./bin/enphase-cli production inverters

# Get raw JSON output
./bin/enphase-cli production --raw
./bin/enphase-cli production inverters --raw

# Use verbose output
./bin/enphase-cli production inverters --verbose
```

### IVP Endpoints
```bash
# Get meter details
./bin/enphase-cli ivp meters

# Get detailed meter readings
./bin/enphase-cli ivp readings

# Get consumption data
./bin/enphase-cli ivp consumption

# Get grid readings
./bin/enphase-cli ivp grid

# Get live system status
./bin/enphase-cli ivp livedata

# Get device list
./bin/enphase-cli ivp ensemble

# Get PDM energy data
./bin/enphase-cli ivp pdm

# Get raw JSON output for any endpoint
./bin/enphase-cli ivp meters --raw
./bin/enphase-cli ivp livedata --raw
./bin/enphase-cli ivp consumption --raw
```

### Output Formats
```bash
# Summary format (default) - Human-readable structured output
./bin/enphase-cli production
./bin/enphase-cli ivp meters

# Raw format - Original JSON response from API
./bin/enphase-cli production --raw
./bin/enphase-cli ivp meters --raw

# Verbose mode - Show connection details and debug info
./bin/enphase-cli production --verbose
./bin/enphase-cli ivp livedata --verbose --raw
```

## 📡 Direct API Usage Examples

### cURL Examples
```bash
# Get system information (no auth)
curl -k https://envoy.local/info

# Get production data
curl -k -H "Authorization: Bearer $TOKEN" \
  https://envoy.local/api/v1/production | jq

# Get inverter details
curl -k -H "Authorization: Bearer $TOKEN" \
  https://envoy.local/api/v1/production/inverters | jq

# Get grid readings
curl -k -H "Authorization: Bearer $TOKEN" \
  https://envoy.local/ivp/meters/gridReading | jq
```

### Go Example
```go
package main

import (
    "crypto/tls"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type ProductionData struct {
    WattsNow          int `json:"wattsNow"`
    WattHoursToday    int `json:"wattHoursToday"`
    WattHoursSevenDays int `json:"wattHoursSevenDays"`
    WattHoursLifetime int `json:"wattHoursLifetime"`
}

func main() {
    gatewayIP := "envoy.local"
    token := "your_bearer_token_here"
    
    // Create HTTP client that ignores SSL certificates
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    }
    
    // Create request
    url := fmt.Sprintf("https://%s/api/v1/production", gatewayIP)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        panic(err)
    }
    
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
    req.Header.Set("Accept", "application/json")
    
    // Make request
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == 200 {
        body, _ := io.ReadAll(resp.Body)
        var data ProductionData
        json.Unmarshal(body, &data)
        
        fmt.Printf("Current production: %d W\n", data.WattsNow)
        fmt.Printf("Today's production: %d Wh\n", data.WattHoursToday)
    } else {
        fmt.Printf("Error: %d\n", resp.StatusCode)
    }
}
```

## 🔧 Hardware Requirements

### Supported IQ Gateway Models
- IQ Gateway (ENV-IQ-AM1-240, ENV2-IQ-AM1-240)
- IQ Gateway Commercial (ENV-IQ-AM3-3P, ENV2-IQC2-AM3-3P)
- IQ Gateway Metered (ENV-S-WM-230, ENV-S-EM-230)
- IQ Gateway Standard (ENV-S-WB-230)
- IQ Gateway M (ENV-S-AM1-230-60)
- Envoy S Standard (ENV-S-AB-120-A)
- Envoy S Metered NA (ENV-S-AM1-120)

### Firmware Requirements
- **Minimum**: IQ Gateway software version 7.0.x or higher
- **Token Support**: Required for local API access

### Optional Hardware (Enables Additional Endpoints)
- **Production CT**: Required for `/api/v1/production` endpoint
- **Consumption CT**: Required for consumption reporting
- **Storage Devices**: Required for battery/storage data
- **Revenue Grade Meters**: For enhanced metering data

## 🎫 Token Management

### Automatic Token Refresh (Shell Script)
```bash
#!/bin/bash
# token_refresh.sh

USERNAME="your_email@example.com"
PASSWORD="your_password"
ENVOY_SERIAL="your_envoy_serial"

get_token() {
    session_id=$(curl -s -X POST http://enlighten.enphaseenergy.com/login/login.json? \
        -F "user[email]=$USERNAME" \
        -F "user[password]=$PASSWORD" | jq -r ".session_id")
    
    web_token=$(curl -s -X POST http://entrez.enphaseenergy.com/tokens \
        -H "Content-Type: application/json" \
        -d "{\"session_id\": \"$session_id\", \"serial_num\": \"$ENVOY_SERIAL\", \"username\": \"$USERNAME\"}")
    
    echo $web_token | jq -r '.token'
}

# Usage
TOKEN=$(get_token)
echo "New token: $TOKEN"
```

### Token Validation
```bash
# Test token validity
curl -k -H "Authorization: Bearer $TOKEN" \
  https://envoy.local/api/v1/production/inverters

# Check response:
# 200 = Valid token
# 401 = Invalid/expired token
```

## ⏱️ Data Update Frequencies

| Data Type | Update Frequency | Endpoints |
|-----------|------------------|-----------|
| **System Info** | Static | `/info` |
| **Device Status** | Real-time | `/ivp/ensemble/device_list` |
| **Meter Readings** | 5 minutes | `/ivp/meters/readings`, `/ivp/meters/reports/consumption` |
| **Production Data** | 5 minutes | `/api/v1/production`, `/api/v1/production/inverters` |
| **Live Data** | Real-time | `/ivp/livedata/status` |
| **Grid Readings** | Real-time | `/ivp/meters/gridReading` |

## 🛠️ Development

### Project Structure
```
enphase/
├── api/
│   └── v1/
│       └── openapi.yaml       # OpenAPI specification
├── cmd/                       # CLI commands
│   ├── root.go               # Root command
│   ├── auth.go               # Authentication commands
│   ├── info.go               # System info command
│   ├── api.go                # Production commands
│   └── ivp.go                # IVP commands
├── internal/                  # Internal packages
│   ├── api/                  # API client
│   └── ivp/                  # IVP client
├── bin/                       # Built binaries
├── Makefile                   # Build targets
├── go.mod                     # Go module
├── main.go                    # Entry point
├── enphase_local.pdf          # Source documentation
└── README.md                  # This file
```

### Building the CLI
```bash
# Build for Linux (default)
make build

# Clean build artifacts
make clean

# Run the application
make run

# Show build help
make help
```

### Code Generation
```bash
# Generate Go client from OpenAPI spec
openapi-generator generate -i api/v1/openapi.yaml -g go -o clients/go

# Generate other clients if needed
openapi-generator generate -i api/v1/openapi.yaml -g <language> -o clients/<language>
```

### Documentation Generation
```bash
# Generate HTML documentation
swagger-codegen generate -i api/v1/openapi.yaml -l html2 -o docs/
```

## 🔍 Troubleshooting

### Common Issues

**Certificate Errors**
```bash
# Use -k flag to ignore self-signed certificates
curl -k https://envoy.local/info
```

**Network Discovery**
```bash
# Find IQ Gateway on network
nmap -sn 192.168.1.0/24 | grep -i envoy
# or
arp -a | grep -i envoy
```

**Token Errors**
- Verify token hasn't expired
- Check if installer vs. owner credentials
- Ensure proper Bearer format: `Authorization: Bearer <token>`

**Connection Timeouts**
- Verify IQ Gateway is on same network
- Check firewall settings
- Ensure IQ Gateway firmware is 7.0.x+

### Debug Mode
```bash
# Verbose curl output
curl -k -v -H "Authorization: Bearer $TOKEN" \
  https://envoy.local/api/v1/production
```

## 📈 Monitoring Setup

### Integration Examples

The CLI tool can be easily integrated into scripts and automation:

```bash
#!/bin/bash
# Simple monitoring script
CURRENT_POWER=$(./bin/enphase-cli production | jq '.wattsNow')
echo "Current solar production: ${CURRENT_POWER}W"

# Log to file
echo "$(date): ${CURRENT_POWER}W" >> production.log
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

**Special Thanks**: This project was developed in collaboration with Claude (Anthropic's AI assistant), who contributed to the code implementation.

### Development Setup
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Update documentation
5. Submit a pull request

### Reporting Issues
- Check existing issues first
- Provide IQ Gateway model and firmware version
- Include relevant error messages
- Share sanitized configuration (remove tokens/credentials)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

This project is based on official Enphase Energy documentation. Please refer to Enphase's terms of service and API usage guidelines.

## 🔗 Related Resources

- [Enphase Developer Portal](https://developer.enphase.com/)
- [Enphase Enlighten Platform](https://enlighten.enphaseenergy.com/)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger Editor](https://editor.swagger.io/)

---

**⚠️ Security Notice**: Always keep your Bearer tokens secure and never commit them to version control. Tokens provide full access to your IQ Gateway's local API.

For questions or support, please open an issue in this repository.
