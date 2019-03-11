package main

func main() {
	// parse flag
	checkOptions()

	if *showVersion {
		printVersion()
	}

	_, err := getStartPages()
	if err != nil {
		exit(1, err)
	}

}
