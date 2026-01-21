package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DryRun      bool
	ApiSubKey   string
	LimiterHost string
	CaCertPath  string
}

func LoadConf() (*Config, error) {
	dryRun := os.Getenv("DRY_RUN")
	dry, err := strconv.ParseBool(dryRun)
	if err != nil {
		return nil, fmt.Errorf("%s is not set", "DRY_RUN")
	}

	apiSubIs := os.Getenv("API_SUBSCRIPTION_KEY")
	if apiSubIs == "" && !dry {
		return nil, fmt.Errorf("%s is not set", "API_SUBSCRIPTION_KEY")
	}

	limiterHost := os.Getenv("LIMITER_HOST")
	if limiterHost == "" && !dry {
		return nil, fmt.Errorf("%s is not set", "LIMITER_HOST")
	}

	caCertPath := os.Getenv("CA_CERT_PATH")
	if caCertPath == "" && !dry {
		return nil, fmt.Errorf("%s is not set", "CA_CERT_PATH")
	}

	return &Config{
		DryRun:      dry,
		ApiSubKey:   apiSubIs,
		LimiterHost: limiterHost,
		CaCertPath:  caCertPath,
	}, nil
}
