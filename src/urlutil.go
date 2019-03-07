package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
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

// discoverHREFURLs discovers HREF urls for fetching
func discoverHREFURLs(s string) (discovered []string) {
	matches := reHREF.FindAllStringSubmatch(s, -1)
	for _, m := range matches {

		url := strings.TrimSpace(m[1])

		if url != "" {
			discovered = append(discovered, m[1])
		}
	}

	return
}

// resolveURL resolve URL based on its parent URL
func resolveURL(parent, child string, unescape bool) (string, error) {
	if unescape {
		parent = html.UnescapeString(parent)
		child = html.UnescapeString(child)
	}

	parsedChild, err := url.Parse(child)

	if err != nil {
		return child, fmt.Errorf("Child -> %v", err)
	}

	if strings.HasPrefix(parsedChild.Host, "http") {
		return parsedChild.String(), nil
	}

	parsedParent, err := url.Parse(parent)
	if err != nil {
		return parsedChild.String(), fmt.Errorf("Parent -> %v",
			err,
		)
	}

	if parsedChild.Host == "" {
		parsedChild.Host = parsedParent.Host
	}

	if parsedChild.Scheme == "" {
		parsedChild.Scheme = parsedParent.Scheme
	}

	if parsedChild.Host != parsedParent.Host {
		return parsedChild.String(), nil
	}

	parsedChild.Path = resolvePath(parsedParent.Path, parsedChild.Path)

	return parsedChild.String(), nil

}

// resolvePath resolves paths to their relative main page path
func resolvePath(page, resource string) string {

	if strings.HasPrefix(resource, "/") {
		return resource
	}

	resource = path.Join(page, resource)
	return path.Clean(resource)

}
