package main

import (
	"akitasoftware.com/demo-client/datasource"
	"github.com/akitasoftware/akita-libs/analytics"
	"github.com/akitasoftware/go-utils/optionals"
	"github.com/golang/glog"
	"math"
	"math/rand"
	"time"
)

type App struct {
	Config *Configuration
	// The demo server used to generate traffic.
	DemoServer datasource.DemoServer
	// The client used to communicate with the Akita backend.
	AkitaClient datasource.AkitaClient
	// Client used to send analytics events.
	AnalyticsClient optionals.Optional[analytics.Client]
}

func NewApp(config *Configuration, server datasource.DemoServer) *App {
	return &App{
		Config:          config,
		DemoServer:      server,
		AkitaClient:     config.AkitaClient,
		AnalyticsClient: config.Analytics,
	}
}

// SendEvent sends an analytics event to the Akita backend.
func (a App) SendEvent(name string, properties map[string]any) {
	// Add the platform to the properties.
	properties["platform"] = a.Config.Platform

	analyticsClient, ok := a.AnalyticsClient.Get()
	if !ok {
		glog.Warning("analytics client not initialized")
		return
	}

	email, err := a.AkitaClient.GetUserEmail(a.Config.Credentials.APIKey, a.Config.Credentials.APISecret)
	if err != nil {
		glog.Errorf("failed to get user email: %v", err)
		return
	}

	if err := analyticsClient.Track(email, name, properties); err != nil {
		glog.Errorf("failed to send analytics event: %v", err)
	}
}

func (a App) HandleDemoTasks() {
	requestInterval := time.Second

	// Create a ticker that fires every second.
	ticker := time.NewTicker(requestInterval)

	// Keep track of the request count.
	requestCount := 0

	for {
		select {
		case <-ticker.C:
			go func() {
				// Send a request to the demo server.
				a.sendMockTraffic()
			}()

			// Increase the request count
			requestCount++

			// Update the request interval
			requestInterval = updateInterval(requestCount)

			// Reset the ticker with the new interval
			ticker.Reset(requestInterval)
		}
	}
}

// updateInterval calculates the new interval based on the request count
func updateInterval(requestCount int) time.Duration {
	// Define a scaling factor to control the rate of interval increase
	scalingFactor := 0.5

	// Calculate the new interval using a logarithmic function
	newInterval := time.Duration(scalingFactor*math.Log10(float64(requestCount+1))) * time.Second

	// Limit the interval to a maximum of 10 seconds
	maxInterval := 10 * time.Second
	if newInterval > maxInterval {
		newInterval = maxInterval
	}

	return newInterval
}

func (a App) sendMockTraffic() {
	handleErr := func(apiName string, err error) {
		glog.Errorf("failed to send demo request to api '%s': %v", apiName, err)
	}

	// To showcase response count metric, we should attempt to send request disproportionately
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomNumber := r.Intn(100)
	if randomNumber < 67 {
		err := a.DemoServer.GetBreed()
		handleErr("GetBreed", err)
	} else {
		err := a.DemoServer.PostTrick()
		handleErr("PostTrick", err)
	}
}
