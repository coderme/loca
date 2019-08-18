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

	stack := filo.NewStringStack()

	for _, page := range pages {
		stack.Push(page)
	}

	concurrent := make(chan struct{}, *concurrency)

	for {

		if stack.Len() == 0 {
			break
		}

		url := stack.Pop()

		if url == "" {
			continue
		}

		concurrent <- struct{}{}

		go func() {
			defer func() { <-concurrent }()
			// process URL here
			//			resp

		}()
	}

}
