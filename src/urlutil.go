package main

import (
	"net/http"
)

var (
	client = &http.Client{}
)

func buildRequest(u, ua string) (*http.Request, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", ua)

	return req, err

}

func fetch(u string) (*http.Response, error) {
	req, err := buildRequest(u, *userAgent)

	if err != nil {
		return nil, err
	}

	return client.Do(req)

}
