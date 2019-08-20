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

func getDir(u string) (string, error) {
	dirCSS := path.Join(*dirAssets, "css")
	dirJS := path.Join(*dirAssets, "js")
	dirFonts := path.Join(*dirAssets, "fonts")
	dirImages := path.Join(*dirAssets, "images")
	dirSVG := path.Join(*dirAssets, "svg")
	// compound paths
	dirPNG := path.Join(dirImages, "png")
	dirJPG := path.Join(dirImages, "jpg")
	// media
	dirVideos := path.Join(*dirMedia, "videos")
	dirMP4 := path.Join(dirVideos, "mp4")
	dirWMV := path.Join(dirVideos, "wmv")
	dirMOV := path.Join(dirVideos, "mov")
	dirAVI := path.Join(dirVideos, "avi")
	dirWEBP := path.Join(dirVideos, "webp")
	dirFLV := path.Join(dirVideos, "flv")
	dirOGV := path.Join(dirVideos, "ogv")
	dirOGX := path.Join(dirVideos, "ogx")
	dirMKV := path.Join(dirVideos, "mkv")

	// audio
	dirAudio := path.Join(*dirMedia, "audios")
	dirAAC := path.Join(dirAudio, "aac")
	dirMP3 := path.Join(dirAudio, "mp3")
	dirMP2 := path.Join(dirAudio, "mp2")
	dirOGG := path.Join(dirVideos, "ogg")
	dirOGA := path.Join(dirVideos, "oga")
	dirM4A := path.Join(dirVideos, "m4a")
	dirWMA := path.Join(dirVideos, "wma")

	mimes := map[string]string{
		// Special cases
		// fonts
		"fonts.googleapis.com/css": dirCSS,
		// Material icons
		"fonts.googleapis.com/icon": dirCSS,
		".css":                      dirCSS,
		".js":                       dirJS,
		".woff2":                    dirFonts,
		".woff":                     dirFonts,
		".oet":                      dirFonts,

		// images
		".svg":  dirSVG,
		".png":  dirPNG,
		".jpg":  dirJPG,
		".jpeg": dirJPG,

		// videos
		".mp4":  dirMP4,
		".mkv":  dirMKV,
		".wmv":  dirWMV,
		".mov":  dirMOV,
		".avi":  dirAVI,
		".webp": dirWEBP,
		".flv":  dirFLV,
		".ogv":  dirOGV,
		".ogx":  dirOGX,
		".m4v":  dirMP4,

		// audio
		".mp2": dirMP2,
		".mp3": dirMP3,
		".m4a": dirM4A,
		".aac": dirAAC,
		".ogg": dirOGG,
		".oga": dirOGA,
		".wma": dirWMA,
	}

	lowered := strings.ToLower(u)

	for i, v := range mimes {
		if strings.Contains(lowered, i) {
			return v, nil
		}
	}

	return guessExt(lowered)

}

// guessExt ties hard to guess extention of a file
func guessExt(u string) (string, error) {

	if strings.Contains(u, ".") {

		parsed, err := url.Parse(u)
		if err != nil {
			return "", err
		}

		return strings.Trim(
			path.Ext(parsed.Path), ".",
		), nil
	}

	return "html", nil

}

// prettyName makes a file name from URL
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
	// err is ignored since it has been dealt with it
	// already
	dir, _ := getDir(u)

	cleaned = append(cleaned,
		cleanedPath(dir),
	)

	// add *cleaned parts to the clean
	cleaned = append(cleaned, parts...)

	return strings.Join(cleaned, "/")

}

// pathParts gets directory and basename from fileName
func dirName(f string) (dir string) {

	parts := strings.Split(f, "/")

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

// rewriteHosts rewrites all href and src URL to 0.0.0.0
// to prevent wasteful overloading of webpages' resources
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

				if !isOfflineHost(parsed.Host) {
					return match0
				}

				return m[1] +
					parsed.Scheme + "://" +
					"0.0.0.0/" + parsed.Host +
					parsed.RequestURI() +
					m[3]

			})
	}

	return data

}

// isOfflineHost lookup host against
func isOfflineHost(host string) bool {
	for _, h := range hosts {
		if strings.Contains(host, h) {
			return true
		}

	}
	return false
}

// mayFetchURL checks whether URL is allowed to be fetched or not
func mayFetchURL(u string) (bool, error) {

	parsed, err := url.Parse(u)

	if err != nil {
		return false, err
	}

	if !*downloadURLsWithQueryString &&
		strings.Contains(u, "?") {
		return false, nil
	}

	if skippableURL(u) {
		return false, nil
	}

	if skippableHost(parsed.Host) {
		return false, nil
	}

	if !allowedURL(u) {
		return false, nil
	}

	if !allowedHost(parsed.Host) {
		return false, nil
	}

	// say yes
	return true, nil
}

// mayFetchContent checks whether Content is allowed to be fetched or not
func mayFetchContent(u string) (bool, error) {

	parsed, err := url.Parse(u)

	if err != nil {
		return false, err
	}

	if !*downloadURLsWithQueryString &&
		strings.Contains(u, "?") {
		return false, nil
	}

	if skippableURL(u) {
		return false, nil
	}

	if skippableHost(parsed.Host) {
		return false, nil
	}

	if !allowedURL(u) {
		return false, nil
	}

	if !allowedHost(parsed.Host) {
		return false, nil
	}

	// say yes
	return true, nil
}

// skippableHost checks if Host may be skipped
func skippableHost(host string) bool {

	hostsSkipped := strings.Split(*skippedHosts, ",")

	for _, h := range hostsSkipped {
		h = strings.TrimSpace(h)

		if h == "" {
			continue
		}

		if host == h {
			return true
		}
	}

	return false
}

// skippableURL checks if URL may be skipped
func skippableURL(url string) bool {

	URLsSkipped := strings.Split(*skippedURLs, ",")

	for _, v := range URLsSkipped {
		v = strings.TrimSpace(v)

		if v == "" {
			continue
		}

		if strings.Contains(url, v) {
			return true
		}
	}

	return false
}

func allowedHost(host string) bool {

	hostsOnly := strings.Split(*onlyHosts, ",")

	var cleanedHostsOnly []string

	for _, h := range hostsOnly {
		h = strings.TrimSpace(h)

		if h != "" {
			cleanedHostsOnly = append(cleanedHostsOnly, h)
		}
	}

	if len(cleanedHostsOnly) > 0 {

		matchHost := false

		for _, h := range cleanedHostsOnly {
			if host == h {
				matchHost = true
				break
			}
		}

		if !matchHost {
			return false
		}
	}

	return true
}

func allowedURL(url string) bool {

	URLsOnly := strings.Split(*onlyURLs, ",")

	var cleanedURLsOnly []string

	for _, v := range URLsOnly {
		v = strings.TrimSpace(v)

		if v != "" {
			cleanedURLsOnly = append(cleanedURLsOnly, v)
		}
	}

	if len(cleanedURLsOnly) > 0 {

		matchURL := false

		for _, v := range cleanedURLsOnly {
			if strings.Contains(url, v) {
				matchURL = true
				break
			}
		}

		if !matchURL {
			return false
		}
	}

	return true
}

// fetchToFile fetch URL and save it to local file
func fetchToFile(u string) error {

	parsed, err := parseURL(u)

	if err != nil {
		return err
	}

	// check URL structure
	// if it allowed to be fetched
	willFetch, err := mayFetchURL(u)
	if err != nil {
		return fmt.Errorf("Err: isAllowedURL(%s) -> err -> %v",
			u,
			err,
		)
	}

	if !willFetch {
		return fmt.Errorf("Err: isAllowedURL(%s) -> NotAllowed",
			u,
		)
	}

	resp, err := fetch(parsed, *delay)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// check URL fetched *content*
	// if it allowed to be stored

	// cool, it seems we plan to save it
	// lets give it a cool name
	name := prettyName(u)

	err = saveFile(resp, name)
	if err != nil {
		return err
	}

	// pretty := prettyURL(u)

	return nil
}

// parseHosts parses data for hosts
func parseHosts(d []byte) []string {
	var hosts []string
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
		hosts = append(hosts, v)
	}

	return hosts

}

// cleanedPath cleans path parts from non-safe
// letters, like space, semicolons
func cleanedPath(p string) string {

	parts := strings.Split(p, "/")
	var cleaned []string

	for _, v := range parts {

		v = prettyURL(v)

		if v == "" {
			continue
		}

		cleaned = append(cleaned, v)
	}

	return strings.Join(cleaned, "/")

}
