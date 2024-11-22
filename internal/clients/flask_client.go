package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type FlaskClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewFlaskClient(baseURL string) *FlaskClient {
	// Ensure the base URL has a scheme
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "http://" + baseURL
	}
	return &FlaskClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// RecognizeFaceWithImage sends image data to the Flask service and returns the recognized face.
func (fc *FlaskClient) RecognizeFaceWithImage(imageData []byte) (string, error) {
	url := fc.BaseURL + "/recognize"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(imageData))
	if err != nil {
		slog.Error("failed to create request", "error", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := fc.HTTPClient.Do(req)
	if err != nil {
		slog.Error("failed to call Flask endpoint", "error", err, "url", url)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Flask service returned error", "status", resp.Status)
		return "", fmt.Errorf("failed to recognize face: %s", resp.Status)
	}

	var result struct {
		RecognizedFace string `json:"recognized_face"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("failed to decode response", "error", err)
		return "", err
	}

	return result.RecognizedFace, nil
}

// RecognizeFace makes a GET request to the Flask service and retrieves the hash.
func (fc *FlaskClient) RecognizeFace() (string, error) {
	url := fc.BaseURL + "/recognize"

	resp, err := fc.HTTPClient.Get(url)
	if err != nil {
		slog.Error("failed to call Flask endpoint", "error", err, "url", url)
		return "", err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Flask service returned error", "status", resp.Status)
		return "", fmt.Errorf("service returned error: %s", resp.Status)
	}

	var result struct {
		Hash  string `json:"hash"`
		Error string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("failed to decode response", "error", err)
		return "", errors.New("could not decode response")
	}

	if result.Hash == "" || result.Error != "" {
		slog.Error("error in response data", "error", result.Error)
		return "", errors.New(result.Error)
	}

	return result.Hash, nil
}

func (fc *FlaskClient) RegisterFace() (string, error) {
	url := fc.BaseURL + "/register"
	resp, err := fc.HTTPClient.Get(url)
	if err != nil {
		slog.Error("failed to call Flask endpoint", "error", err, "url", url)
		return "", err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if resp.StatusCode != http.StatusOK {
		slog.Error("Flask service returned error", "status", resp.Status)
		return "", fmt.Errorf("service returned error: %s", resp.Status)
	}

	var result struct {
		Hash  string `json:"hash"`
		Error string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("failed to decode response", "error", err)
		return "", errors.New("could not decode response")
	}

	if result.Hash == "" || result.Error != "" {
		slog.Error("error in response data", "error", result.Error)
		return "", errors.New(result.Error)
	}

	return result.Hash, nil
}
