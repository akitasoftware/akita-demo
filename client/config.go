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
	} `yaml:"akita"`
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
		glog.Infof("Enabling analytics...")
		client, err := analytics.NewClient(rawConfig.Analytics.Config)
		if err != nil {
			// If we fail to create the analytics client, we don't want to fail the entire demo.
			// Instead, we just log the error and continue without analytics.
			glog.Errorf("Failed to create analytics client: %v", err)
		}

		glog.Infof("Analytics client created successfully")
		analyticsClient = optionals.Some(client)
	} else {
		glog.Infof("Analytics have been disabled")
	}

	return &Configuration{
		AkitaClient: akitaClient,
		Analytics:   analyticsClient,
		Credentials: struct {
			APIKey    string
			APISecret string
		}{
			// The API key and secret are set by the Docker Compose file.
			APIKey:    os.Getenv("AKITA_API_KEY_ID"),
			APISecret: os.Getenv("AKITA_API_KEY_SECRET"),
		},
		// Get the target platform declared at build time.
		Platform: os.Getenv("TARGETPLATFORM"),
	}, nil
}
