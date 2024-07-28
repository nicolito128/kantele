/*
rest includes the necessary implementations
to make HTTP requests to the Discord API.
*/
package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	Version    string = "v10"
	Domain     string = "https://discord.com"
	DiscordAPI string = Domain + "/api/" + Version
)

const (
	github = "github.com/nicolito128/kantele"
)

type Rest struct {
	botToken   string
	httpClient *http.Client
}

func New(token string) *Rest {
	return &Rest{
		botToken:   token,
		httpClient: &http.Client{},
	}
}

// Do sends an HTTP request using an internal client,
// attaching the necesary authorization headers,
// then returns an HTTP response and an error.
func (r *Rest) Do(method, endpoint string, data any) (*http.Response, error) {
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}
	url := fmt.Sprintf("%s%s", DiscordAPI, endpoint)

	// Encoding data to request body
	var out []byte
	var err error

	if data != nil {
		out, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("encoding request data error: %w", err)
		}
	}

	// Making the request to the discord api
	req, err := http.NewRequest(method, url, bytes.NewReader(out))
	if err != nil {
		return nil, fmt.Errorf("doing request error: %w", err)
	}

	// Adding headers
	req.Header.Add("User-Agent", fmt.Sprintf("DiscordBot (%s)", github))
	req.Header.Add("Authorization", fmt.Sprintf("Bot %s", r.botToken))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http send error: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"invalid response: CODE: %d METHOD: %s PATH: %s URL: %s",
			res.StatusCode,
			method,
			res.Request.URL.Path,
			url,
		)
	}

	return res, nil
}

// Get sends an HTTP GET request using an internal client,
// then returns an HTTP response and an error.
func (r *Rest) Get(endpoint string) (*http.Response, error) {
	return r.Do(http.MethodGet, endpoint, nil)
}

// Post sends an HTTP POST request using an internal client,
// then returns an HTTP response and an error.
func (r *Rest) Post(endpoint string, data any) (*http.Response, error) {
	return r.Do(http.MethodPost, endpoint, data)
}

// Patch sends an HTTP PATCH request using an internal client,
// then returns an HTTP response and an error.
func (r *Rest) Patch(endpoint string, data any) (*http.Response, error) {
	return r.Do(http.MethodPatch, endpoint, data)
}

// Delete sends an HTTP DELETE request using an internal client,
// then returns an HTTP response and an error.
func (r *Rest) Delete(endpoint string) (*http.Response, error) {
	return r.Do(http.MethodDelete, endpoint, nil)
}
