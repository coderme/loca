package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	version                   = 0.1
	defaultUserAgent          = `Mozilla/5.0 Loca (%v) (https://github/codermeorg/loca)`
	defaultDelayBeforeRequest = 3 * time.Second
	defaultConcurrency        = 1
	retryDefaultCount         = 3
	// default
	defaultDir          = `./`
	defaultDirAssets    = `assets`
	defaultDirMedia     = `media`
	defaultDirArchieve  = `archieves`
	defaultDirDocs      = `docs`
	defaultDirPages     = "pages"
	defaultDirUnsorted  = "unsorted"
	defaultLang         = "en"
	defaultSkippedHosts = "youtube.com,vimeo.com,github.com,bitbucket.com"
	defaultSkippedURLs  = ""
	defaultOnlyHosts    = ""
	defaultOnlyURLs     = ""
	tempFilePrefix      = "loca-"
	defaultOfflineList  = "https://raw.githubusercontent.com/codermeorg/loca-offline/master/src/hosts"
)

var (
	concurrency = flag.Int("c", defaultConcurrency, "Concurrency level for parallel URL fetching")
	delay       = flag.Duration("delay", defaultDelayBeforeRequest, "")
	langs       = flag.String("langs", defaultLang, "Download pages with these langs: CSV or single language code")

	downloadMedia               = flag.Bool("dl-media", false, "Download videos and audio files")
	downloadArchieve            = flag.Bool("dl-archieve", false, "Download archieve files")
	downloadURLsWithQueryString = flag.Bool("dl-query", false, "Download URLs those with query string like https://example.com/?action=msglist&order=desc")
	ascend                      = flag.Bool("ascend", false, "Ascend to host root when fetching")
	retry                       = flag.Bool("retry", false, "Retry fetching directly if fetch failed")
	retryCount                  = flag.Int("retry-max-count", retryDefaultCount, "Retry fetching this times before giving up")
	verbose                     = flag.Bool("verbose", false, "Be more verbose")

	dir          = flag.String("dir", defaultDir, "Dirctory root where to store all downloaded files.")
	dirAssets    = flag.String("dir-assets", defaultDirAssets, "Dirctory where to store assets files.")
	dirMedia     = flag.String("dir-media", defaultDirMedia, "Dirctory where to store videos and audios files.")
	dirDocs      = flag.String("dir-Docs", defaultDirDocs, "Dirctory where to store doc, epub and pdf.")
	dirPages     = flag.String("dir-pages", defaultDirPages, "Dirctory where to store HTML pages.")
	dirArchieves = flag.String("dir-archieves", defaultDirArchieve, "Directory where to store archieve files")
	dirUnsorted  = flag.String("dir-unsorted", defaultDirUnsorted, "Dirctory where to store Unsorted files.")

	skippedHosts = flag.String("skipped-hosts", defaultSkippedHosts, "CSV, skip fetching any host that contains any of these values.")
	skippedURLs  = flag.String("skipped-urls", defaultSkippedURLs, "CSV, skip fetching any url that contains any of these values.")

	onlyHosts = flag.String("only-hosts", defaultOnlyHosts, "CSV, Fetch only hosts that contain any of these values.")
	onlyURLs  = flag.String("only-urls", defaultOnlyURLs, "CSV, Fetch only URLs that contain any of these values.")

	userAgent       = flag.String("user-agent", defaultUserAgent, "UserAgent of the client")
	keepMeta        = flag.Bool("keep-meta", false, "Keep original <meta> tags")
	offlineDisabled = flag.Bool("offline-disabled", false, "Disable rewriting hosts for offline browsing")

	offlineHosts = flag.String("offline-list", defaultOfflineList, "List of websites to be rewriting for offline browsing")

	showVersion = flag.Bool("v", false, "Print version")

	// global vars
	hosts []string
)

func init() {
	if os.Geteuid() == 0 || os.Getegid() == 0 {
		exit(2, "Don't Run as ROOT")
	}

}

// parseOptions checks the command args for validity
func parseOptions() error {
	flag.Parse()

	if *concurrency <= 0 {
		*concurrency = defaultConcurrency
	}

	if *retryCount < 0 {
		*retryCount = 0
	}

	args := flag.Args()

	if len(args) == 0 {
		return fmt.Errorf("Missing URL(s) to work on")
	}

	if !*offlineDisabled {
		var err error
		hosts, err = cacheHosts(*offlineHosts)

		if err != nil {
			return fmt.Errorf("cacheHosts()-> %v", err)
		}

	}

	return nil

}

//getStartPages gets start page
func getStartPages() ([]string, error) {

	var pages []string
	for _, u := range flag.Args() {
		parsed, err := parseURL(u)
		if err != nil && *verbose {
			log.Println("Error: parseing", u, "->", err)
			continue
		}
		pages = append(pages, parsed)
	}

	if len(pages) == 0 {
		return nil, fmt.Errorf("Error: no valid URL provided")
	}
	return pages, nil

}

// printVersion prints loca version
// then exit
func printVersion() {
	fmt.Fprintf(os.Stderr, `%v %v (c)(%v)`,
		os.Args[0],
		version,
		time.Now().Year(),
	)
	exit(1)
}

// cacheHosts graps hosts from remote URL or local file
func cacheHosts(u string) ([]string, error) {

	// local file
	if strings.HasPrefix(u, "/") {
		data, err := ioutil.ReadFile(u)
		if err != nil {
			return nil, err
		}

		return parseHosts(data), nil
	}

	resp, err := fetch(u, time.Nanosecond)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err

	}

	return parseHosts(data), nil

}
