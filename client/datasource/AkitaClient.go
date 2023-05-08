package datasource

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

// AkitaClient is the interface for communicating with the Akita backend.
type AkitaClient interface {
	// GetUserEmail returns the email of the user associated with the API key and secret.
	GetUserEmail(APIKey, APISecret string) (string, error)
}

type akitaClientImpl struct {
	baseURL    string
	httpClient *http.Client
}

// Creates a new AkitaClient with the given base URL and HTTP client.
func NewAkitaClient(baseURL string, httpClient *http.Client) AkitaClient {
	return &akitaClientImpl{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (a akitaClientImpl) GetUserEmail(APIKey, APISecret string) (string, error) {
	const path = "/v1/user"

	type payload struct {
		Email string `json:"email"`
	}

	var response payload

	req, err := http.NewRequest(http.MethodGet, a.baseURL+path, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to create request to fetch user email")
	}

	req.SetBasicAuth(APIKey, APISecret)

	resp, err := a.httpClient.Do(req)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response body while fetching user email")
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.Wrapf(
			err,
			"failed to fetch user email. status code: %d. body: %s",
			resp.StatusCode,
			string(body),
		)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", errors.Wrapf(err, "failed to decode response body from user request. body: %s", string(body))
	}

	if response.Email == "" {
		return "", errors.Errorf("response body from user request did not contain email. body: %s", string(body))
	}

	return response.Email, nil
}
