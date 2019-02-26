package main

import (
	"io/ioutil"
	"net/http"
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
