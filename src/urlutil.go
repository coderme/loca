package main

import (
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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

func guessMime(u string) string {

	return ""
}

// prettyName makes a file form URL
func prettyName(u string) string {
	// get rid of the scheme://
	url := strings.SplitN(u, "://", 2)[1]

	// assume url uses / as seprator, then
	// we need to clean it
	parts := strings.Split(url, "/")

	for i, v := range parts {
		parts[i] = prettyURL(v)
	}

	var cleaned []string

	// get dir based on guessed mimetime

	// add *cleaned parts to the clean
	cleaned = append(cleaned, parts...)

	// return ""

}

// pathParts gets directory and basename from fileName
func dirName(fileName string) (dir string) {

	parts := strings.Split(fileName, "/")

	return strings.Join(parts[:len(parts)-1], "/")

}

// saveFile save the resp body HTTP response to name file
func saveFile(resp *http.Response, name string) error {
	defer resp.Body.Close()

	f, err := ioutil.TempFile("", tempFilePrefix)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		defer func() {
			// clean the mess
			err := os.Remove(f.Name())
			if err != nil {
				log.Println(err)
			}
		}()

		return err
	}

	dir := dirName(name)

	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	err = os.Rename(f.Name(), name)

	return err

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
		return validated, fmt.Errorf("invalid URL")
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

	page = path.Dir(page)

	if strings.HasPrefix(resource, "/") {
		return resource
	}

	resource = path.Join(page, resource)
	return path.Clean(resource)

}

// rewriteHosts rewrites all href and frc URL to 0.0.0.0
// to prevent never-ending loading status of webpages
// while browsing offline
func rewriteOfflineURLs(data string) string {
	for _, re := range reFilter {

		data = re.ReplaceAllStringFunc(data,
			func(match0 string) string {
				m := re.FindStringSubmatch(match0)
				parsed, err := url.Parse(m[2])
				if err != nil {
					return m[1] + "#0" + m[3]
				}

				if !isOfflineHost(m[2]) {
					return match0
				}

				return m[1] +
					parsed.Scheme + "://" +
					"0.0.0.0" + parsed.Host +
					parsed.RequestURI() +
					m[3]

			})
	}

	return data

}

// isOfflineHost lookup host against
func isOfflineHost(host string) bool {
	_, ok := hosts[host]
	return ok
}

// fetchToFile fetch URL and save it to local file
// system
func fetchToFile(u string) error {

	parsed, err := parseURL(u)

	if err != nil {
		return err
	}

	// check URL structure
	// if it allowed to be fetched

	resp, err := fetch(parsed, *delay)
	if err != nil {
		return err
	}

	// check URL fetched content
	// if it allowed to be stored

	// cool, it seems we plan to save it
	// lets name it
	name := prettyName(u)

	err = saveFile(resp, name)
	if err != nil {
		return err
	}

	pretty := prettyURL(u)

	return nil
}

// parseHosts parses data for hosts
func parseHosts(d []byte) map[string]bool {
	hosts := make(map[string]bool)
	lines := strings.Split(string(d), "\n")

	for _, v := range lines {
		v = strings.TrimSpace(v)
		// skip comments or empty v
		if strings.HasPrefix(v, "#") ||
			v == "" {
			continue
		}
		// if v has whitespace, host should not
		if strings.Contains(v, " ") {
			parts := strings.Split(v, " ")
			// we pick th first value
			// the rest likely to be a comment
			v = parts[0]
		}
		hosts[v] = true
	}

	return hosts

}
