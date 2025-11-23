package api

import (
	"log"
	"net/http"
	"net/url"

	"tahanraamatut/internal/dotenv"
)

type ReadarrAPI struct {
	dotenv dotenv.Dotenv
	client *http.Client
}

func Connect() ReadarrAPI {
	env := dotenv.GetEnv()
	return ReadarrAPI{
		dotenv: env,
		client: &http.Client{},
	}
}

func (api ReadarrAPI) makeRequest(PathURL *url.URL) (*http.Response, error) {
	req, err := http.NewRequest("GET", PathURL.String(), nil)
	if err != nil {
		log.Fatalf("encountered a fatal error: %s", err)
	}

	log.Println(api.dotenv.ApiKey)

	req.Header.Add("X-Api-Key", api.dotenv.ApiKey)
	req.Header.Add("Accept", "*/*")

	resp, err := api.client.Do(req)

	return resp, err
}

func (api ReadarrAPI) HealthCheck() (*http.Response, error) {
	envVar := api.dotenv

	healthURL := envVar.ApiURL.JoinPath("api", "v1", "health")

	resp, err := api.makeRequest(healthURL)

	return resp, err
}

func (api ReadarrAPI) Search(term string) (*http.Response, error) {
	envVar := api.dotenv

	searchURL := envVar.ApiURL.JoinPath("api", "v1", "search")

	queryValues := searchURL.Query()
	queryValues.Add("term", term)

	searchURL.RawQuery = queryValues.Encode()

	resp, err := api.makeRequest(searchURL)

	return resp, err
}
