package main

import (
	"net/http"
)

var (
	client = &http.Client{}
)

func buildRequest(u, ua string) *http.Request {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", ua)

	return ua, err

}

func fetch(u string) (*http.Response, error) {

}
