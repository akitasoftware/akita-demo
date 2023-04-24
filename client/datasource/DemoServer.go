package datasource

import (
	"bytes"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"net/http"
)

var (
	// List of tricks IDs and their corresponding probability of being selected when sending a request to the mock demo server.
	tricks map[string]float32
)

type (
	DemoServer interface {
		// Send a random breed request to the demo server.
		GetBreed() error
		// Send a random trick request to the demo server.
		PostTrick() error
	}
	demoServerImpl struct {
		baseURL    string
		httpClient http.Client
	}
)

func NewDemoServer(baseURL string) DemoServer {
	return &demoServerImpl{
		baseURL: baseURL,
	}
}

func (d demoServerImpl) GetBreed() error {
	breedID, err := getRandomBreedID()
	if err != nil {
		return fmt.Errorf("failed to get random breed ID: %w", err)
	}

	url := fmt.Sprintf("%s/breeds/%s", d.baseURL, breedID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")

	_, err = d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	return nil
}

func getRandomBreedID() (string, error) {
	breeds := []any{
		// Gives a 20% chance of returning a 404
		"4e7bde8a-92a6-4a4a-a1e9-5547537e90f7",
		"33f9889c-e4aa-4ef4-ba2d-560c1048bc9b",
		"dcd6b113-19a1-41af-8037-84c02951b990",
		"09348399-fb03-4fcc-9a4b-a1eaf796bd75",
		// Gives an 80% chance of returning a 200
		gofakeit.UUID(),
	}

	probabilities := []float32{
		0.05,
		0.05,
		0.05,
		0.05,
		0.80,
	}

	selectedBreed, err := gofakeit.Weighted(breeds, probabilities)
	if err != nil {
		return "", fmt.Errorf("failed to pick a random breed: %w", err)
	}

	return selectedBreed.(string), nil
}

func (d demoServerImpl) PostTrick() error {
	trickID, err := getRandomTrickID()
	if err != nil {
		return fmt.Errorf("failed to get random trick ID: %w", err)
	}

	url := fmt.Sprintf("%s/tricks/%s", d.baseURL, trickID)

	// Prepare a request body with a random owner ID, name, address, and treat count.
	body := fmt.Sprintf(
		`{ owner: { id: "%s", name: "%s", address: "%s" }, treat_count: %d}`,
		gofakeit.UUID(),
		gofakeit.Name(),
		gofakeit.Address().Address,
		gofakeit.Number(1, 10),
	)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	_, err = d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to treat send request: %w", err)
	}

	return nil
}

func getRandomTrickID() (string, error) {
	tricks := []any{
		// Gives a 10% chance to cause a 400
		"bb5a4789-8189-4905-a736-682de6a32375",
		"69d48609-ac34-4d36-bd7f-46f1207ee80e",
		// Gives a 10% chance to cause a 500
		"dc722acb-45e1-4e3e-a926-b186929e6570",
		"f2821a1d-b5f6-4a16-a1ed-b78fce03703d",
		// Gives an 80% chance of returning a 200
		gofakeit.UUID(),
	}

	probabilities := []float32{
		0.05,
		0.05,
		0.05,
		0.05,
		0.8,
	}

	selectedTrick, err := gofakeit.Weighted(tricks, probabilities)
	if err != nil {
		return "", fmt.Errorf("failed to pick a random trick: %w", err)
	}

	return selectedTrick.(string), nil
}
