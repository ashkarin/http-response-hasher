package app_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"http-response-hasher/app"
	"http-response-hasher/hasher"
)

func TestApp(t *testing.T) {
	// Establish a test server
	serverResponse := "Test server response"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, serverResponse)
	})
	ts := httptest.NewServer(handler)

	// Compute expected hash string
	expectedHashStr := hasher.HashToStr(hasher.ComputeHash([]byte(serverResponse)))

	// Start test
	urls := []string{ts.URL}
	results, err := app.ProcessUrls(urls, 1)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	for result := range results {
		if string(result.Input) != ts.URL {
			t.Errorf("Result has a wrong input %v != %s", result.Input, ts.URL)
		}
		if string(result.Output) != expectedHashStr {
			t.Errorf("Result has a wrong output %v != %v", result.Output, expectedHashStr)
		}
		if result.Error != nil {
			t.Errorf("Result has unexpected error: %v", result.Error)
		}
	}
}
