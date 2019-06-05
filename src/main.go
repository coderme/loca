package main

import (
	"github.com/codermeorg/filo"
)

func main() {
	// parse flag
	parseOptions()

	if *showVersion {
		printVersion()
	}

	pages, err := getStartPages()
	if err != nil {
		exit(1, err)
	}

	// init queue
	queue := filo.NewStringStack()

	for _, page := range pages {
		queue.Push(page)
	}

	concurrent := make(chan struct{}, *concurrency)

	for {

		if queue.Len() == 0 {
			break
		}

		url := queue.Pop()

		if url == "" {
			continue
		}

		concurrent <- struct{}{}

		go func() {
			defer <-concurrent
			// process URL here

			

		}()
	}

}
