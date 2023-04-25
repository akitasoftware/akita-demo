package main

import (
	"akitasoftware.com/demo-client/datasource"
	_ "embed"
	"fmt"
	"github.com/akitasoftware/akita-libs/analytics"
	"github.com/akitasoftware/go-utils/optionals"
	"github.com/golang/glog"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
)

type rawConfiguration struct {
	Analytics struct {
		// Configures the analytics client.
		analytics.Config `yaml:",inline"`
		// Whether analytics are enabled.
		Enabled bool `yaml:"enabled"`
	} `yaml:"analytics"`
	Akita struct {
		// Base URL of the Akita API.
		BaseURL string `yaml:"base_url"`
	}
}

type UserCredentials struct {
	// The credentials used to identify the user.
	APIKey    string
	APISecret string
}

type Configuration struct {
	AkitaClient datasource.AkitaClient
	Analytics   optionals.Optional[analytics.Client]
	Credentials UserCredentials
	// The target platform for the Docker image the client is derived from.
	Platform string `yaml:"platform"`
}

func ParseConfiguration(rawData []byte) (*Configuration, error) {
	var rawConfig rawConfiguration
	if err := yaml.Unmarshal(rawData, &rawConfig); err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Create the Akita client.
	akitaClient := datasource.NewAkitaClient(rawConfig.Akita.BaseURL, http.DefaultClient)

	// Create the analytics client.
	var analyticsClient optionals.Optional[analytics.Client]
	if rawConfig.Analytics.Enabled {
		client, err := analytics.NewClient(rawConfig.Analytics.Config)
		if err != nil {
			// If we fail to create the analytics client, we don't want to fail the entire demo.
			// Instead, we just log the error and continue without analytics.
			glog.Errorf("failed to create analytics client: %v", err)
		}

		analyticsClient = optionals.Some(client)
	}

	akitaAPIKey := os.Getenv("AKITA_API_KEY")
	akitaAPISecret := os.Getenv("AKITA_API_SECRET")
	platform := os.Getenv("TARGETPLATFORM")

	return &Configuration{
		AkitaClient: akitaClient,
		Analytics:   analyticsClient,
		Credentials: struct {
			APIKey    string
			APISecret string
		}{
			APIKey:    akitaAPIKey,
			APISecret: akitaAPISecret,
		},
		Platform: platform,
	}, nil
}