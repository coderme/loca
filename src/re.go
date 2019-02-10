package main

import (
	"regexp"
)

var (
	reHREF = regexp.MustCompile(`(?i) href=["']([^<>"']+)['"]`)
	reSRC  = regexp.MustCompile(`(?i) src=["']([^<>"']+)['"]`)

	// integrety and checksum are problematic on localhost
	reWhitespace = regexp.MustCompile(`(?s)\s+`)
)
