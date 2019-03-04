package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	client = &http.Client{}
)

// buildRequest builds a HTTP request and sets a custom User Agent
func buildRequest(u, ua string) (*http.Request, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", ua)

	return req, err

}

// fetch fetches a HTTP resource after the delay
func fetch(u string, delay time.Duration) (*http.Response, error) {
	// delayed fetching
	<-time.After(delay)

	req, err := buildRequest(u, *userAgent)

	if err != nil {
		return nil, err
	}

	return client.Do(req)
}

// readFullResponse reads the full HTTP response as bytes, closes
// HTTP response bdy after finish
func readFullResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)

}

// parseURL parses URL returns valid URL for fetching,
// and any error encountered while parsing
func parseURL(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return u, err
	}

	if parsed.Scheme == "" {
		parsed.Scheme = "http"
	}

	validated := parsed.String()

	if !strings.HasPrefix(validated, "http") {
		return validated, fmt.Errorf("invalid URL supported")
	}

	return validated, err

}

// discoverSrcURLs discovers SRC urls for fetching
func discoverSrcURLs(s string) (discovered []string) {
	matches := reSRC.FindAllStringSubmatch(s, -1)
	for _, m := range matches {

		url := strings.TrimSpace(m[1])

		if url != "" {
			discovered = append(discovered, m[1])
		}
	}

	return
}
