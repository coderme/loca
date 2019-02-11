package main

import (
	"regexp"
)

var (
	reHREF    = regexp.MustCompile(`(?i) href=["']([^<>"']+)['"]`)
	reSRC     = regexp.MustCompile(`(?i) src=["']([^<>"']+)['"]`)
	reUglyURL = regexp.MustCompile(`(?s)[;\+\s"'` + "`" + `%\?!~><\}\{\]\[\\\\:,|\*&^$@]+`)

	// integrety and crossorigin maybe problematic on localhost
	reCrossOrigin = regexp.MustCompile(`(?is)<[^<>]+(\scrossorigin=["'][^"'<>]*["'])`)
	reIntegrity   = regexp.MustCompile(`(?is)<[^<>]+(\sintegrity=["'][^"'<>]*["'])`)
	reWhitespace  = regexp.MustCompile(`(?s)\s+`)
)
