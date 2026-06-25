package config

import (
	"os"
	"strconv"
)

// Config holds the runtime configuration assembled from environment variables
// and global CLI flags. Credentials come from the OKX API-management page; all
// three (key, secret, passphrase) are required to sign private requests.
var Config struct {
	APIKey     string
	APISecret  string
	Passphrase string
	Proxy      string
	BaseURL    string
	Demo       bool
	OutputJSON bool
}

func init() {
	Config.APIKey = os.Getenv("OKX_API_KEY")
	Config.APISecret = os.Getenv("OKX_API_SECRET")
	Config.Passphrase = os.Getenv("OKX_PASSPHRASE")
	Config.Proxy = os.Getenv("OKX_PROXY")
	Config.BaseURL = os.Getenv("OKX_BASE_URL")
	if v := os.Getenv("OKX_DEMO"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			Config.Demo = b
		}
	}
}
