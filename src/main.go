package main

func main() {
	// parse flag
	checkOptions()

	if *showVersion {
		printVersion()
	}

}
