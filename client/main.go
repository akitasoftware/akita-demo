package main

import (
	"akitasoftware.com/demo-client/datasource"
	"context"
	"math"
	"math/rand"
	"net/http"
	"time"
)

const (
	demoServerURL = "http://akita-demo-server:8080"
)

func main() {
	appCtx := context.Background()

	// Create a new demo server
	demoServer := datasource.NewDemoServer(demoServerURL, http.DefaultClient)

	// Start sending requests to the demo server
	handleDemoTasks(appCtx, demoServer)
}

func sendMockTraffic(ctx context.Context, demoServer datasource.DemoServer) {
	handleErr := func(apiName string, err error) {
		if err != nil {
			// TODO: Log Error
		}
	}

	// To showcase response count metric, we should attempt to send request disproportionately
	// TODO: Figure out which isn't deprecated
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

func handleDemoTasks(ctx context.Context, demoServer datasource.DemoServer) {
	requestInterval := time.Second

	// Create a ticker with the initial interval
	ticker := time.NewTicker(requestInterval)

	// Keep track of the request count
	requestCount := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Send a request
			sendMockTraffic(ctx, demoServer)

			// Increase the request count
			requestCount++

			// Update the request interval
			updateInterval(requestCount)

			// Reset the ticker with the new interval
			ticker = time.NewTicker(requestInterval)
		}
	}
}

// updateInterval calculates the new interval between requests based on the current request count.
// It uses a custom base logarithm function to achieve a slow tapering off of request frequency.
func updateInterval(requestCount int) time.Duration {
	// Define the custom base for the logarithm function.
	// A user should reach around 10k requests in 30 minutes.
	base := 1.2

	// Calculate the new interval by taking the custom base logarithm of the request count (plus 1 to avoid log(0)).
	// The result is then multiplied by the duration of 1 second to get the new interval in time.Duration format.
	newInterval := time.Duration(math.Log(float64(requestCount+1))/math.Log(base)) * time.Second

	// Return the updated interval.
	return newInterval
}
