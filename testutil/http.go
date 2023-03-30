package testutil

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Post executes POST requests
func Post(t *testing.T, url, body string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	assert.NoError(t, err)
	for header, headerVal := range headers {
		req.Header.Add(header, headerVal)
	}
	return client.Do(req)
}

// Get executes GET requests
func Get(t *testing.T, url string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest("GET", url, http.NoBody)
	assert.NoError(t, err)
	for header, headerVal := range headers {
		req.Header.Add(header, headerVal)
	}
	return client.Do(req)
}
