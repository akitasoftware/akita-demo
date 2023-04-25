package main

import (
	"akitasoftware.com/demo-client/datasource"
	_ "embed"
	"github.com/akitasoftware/akita-libs/analytics"
	"math"
	"math/rand"
	"net/http"
	"time"
)

//go:embed application.yml
var applicationYML []byte

const (
	demoServerURL = "http://akita-demo-server:8080"
)

func main() {
	analytics.NewClient(analytics.Config{})

	// Create a new demo server
	demoServer := datasource.NewDemoServer(demoServerURL, http.DefaultClient)

	// Start sending requests to the demo server
	handleDemoTasks(demoServer)
}

func sendMockTraffic(demoServer datasource.DemoServer) {
	handleErr := func(apiName string, err error) {
		if err != nil {
			// TODO: Log Error
		}
	}

	// To showcase response count metric, we should attempt to send request disproportionately
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomNumber := r.Intn(100)
	if randomNumber < 67 {
		err := demoServer.GetBreed()
		handleErr("GetBreed", err)
	} else {
		err := demoServer.PostTrick()
		handleErr("PostTrick", err)
	}
}

func handleDemoTasks(demoServer datasource.DemoServer) {
	requestInterval := time.Second

	// Create a ticker with the initial interval
	ticker := time.NewTicker(requestInterval)

	// Keep track of the request count
	requestCount := 0

	for {
		select {
		case <-ticker.C:
			// Send a request
			sendMockTraffic(demoServer)

			// Increase the request count
			requestCount++

			// Update the request interval
			updateInterval(requestCount)

			// Reset the ticker with the new interval
			ticker = time.NewTicker(requestInterval)
		}
	}
}

// updateInterval calculates the new interval based on the request count
func updateInterval(requestCount int) time.Duration {
	// Define a scaling factor to control the rate of interval increase
	scalingFactor := 0.5

	// Calculate the new interval using a logarithmic function
	newInterval := time.Duration(scalingFactor*math.Log10(float64(requestCount+1))) * time.Second

	// Limit the interval to a maximum of 30 seconds
	maxInterval := 30 * time.Second
	if newInterval > maxInterval {
		newInterval = maxInterval
	}

	return newInterval
}
