package main

import (
	"akitasoftware.com/demo-client/datasource"
	"github.com/akitasoftware/akita-libs/analytics"
	"github.com/akitasoftware/go-utils/optionals"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"math/rand"
	"sync"
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

// Create a new app instance with the given configuration and demo server.
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
		glog.Warningf("analytics client not configured")
		return
	}

	email, err := a.AkitaClient.GetUserEmail(a.Config.Credentials.APIKey, a.Config.Credentials.APISecret)
	if err != nil {
		glog.Errorf("failed to get user email: %v", err)
		return
	}

	if err := analyticsClient.Track(email, name, properties); err != nil {
		glog.Errorf("failed to emit analytics event: %v", err)
	}
}

// HandleDemoTasks sends requests to the demo server at a regular interval.
func (a App) HandleDemoTasks() {
	requestInterval := time.Second

	// Create a ticker that fires every second.
	ticker := time.NewTicker(requestInterval)
	// Mutex for keeping track of errors logged by the ticker.
	var rwMutex sync.RWMutex
	var errorCount int

	for {
		select {
		case <-ticker.C:
			go func() {
				// Send a request to the demo server.
				err := a.sendMockTraffic()

				// Send an error event if we've sent less than 5 error events
				// TODO: It would be nice to have the error count be configurable.
				// We could also consider resetting the error count after a certain amount of time.
				rwMutex.RLock()
				if err != nil && errorCount < 5 {
					a.SendEvent(
						"Demo Client Error", map[string]any{
							"error": err.Error(),
						},
					)
					// Increment the error count. This is protected by a mutex because we're in a goroutine.
					rwMutex.RUnlock()
					rwMutex.Lock()
					errorCount++
					rwMutex.Unlock()
				}
			}()
		}
	}
}

// Send a random request to the demo server.
func (a App) sendMockTraffic() error {
	var err error

	// To showcase response count metric, we should attempt to send request disproportionately
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomNumber := r.Intn(100)
	if randomNumber < 10 {
		err = a.DemoServer.GetOwner()
	} else if randomNumber < 67 {
		err = a.DemoServer.GetBreed()
	} else {
		err = a.DemoServer.PostTrick()
	}

	if err != nil {
		return errors.Wrap(err, "failed to send mock traffic")
	}

	return nil
}
