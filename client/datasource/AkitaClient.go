package datasource

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

type AkitaClient interface {
	GetUserEmail(APIKey, APISecret string) (string, error)
}

type akitaClientImpl struct {
	baseURL    string
	httpClient *http.Client
}

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

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode response body")
	}

	return response.Email, nil
}
