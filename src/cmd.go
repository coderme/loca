package main

import (
	"flag"
	"fmt"
	"os"
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
	defaultLang         = "en"
	defaultSkippedHosts = "youtube.com,vimeo.com"
	defaultSkippedURLs  = ""
	defaultOnlyHosts    = ""
	defaultOnlyURLs     = ""
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
	dirArchieves = flag.String("dir-archieves", defaultDirArchieve, "Directory where to store archieve files")

	skippedHosts = flag.String("skipped-hosts", defaultSkippedHosts, "CSV, skip fetching any host that contains any of these values.")
	skippedURLs  = flag.String("skipped-urls", defaultSkippedURLs, "CSV, skip fetching any url that contains any of these values.")

	onlyHosts = flag.String("only-hosts", defaultOnlyHosts, "CSV, Fetch only hosts that contain any of these values.")
	onlyURLs  = flag.String("only-urls", defaultOnlyURLs, "CSV, Fetch only URLs that contain any of these values.")

	userAgent = flag.String("user-agent", defaultUserAgent, "UserAgent of the client")
	keepMeta  = flag.Bool("keep-meta", false, "Keep origin <meta> tags")

	showVersion = flag.Bool("v", false, "Print version")
)

func init() {
	if os.Geteuid() == 0 || os.Getegid() == 0 {
		fmt.Println("Don't Run as ROOT")
		os.Exit(2)
	}

}

// checkOptions checks the command args for validity
func checkOptions() error {
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

	return nil

}

// printVersion prints loca version
// then exit
func printVersion() {
	fmt.Fprintf(os.Stderr, `%v %v (c)(%v)`,
		os.Args[0],
		version,
		time.Now().Year(),
	)
	os.Exit(1)
}
