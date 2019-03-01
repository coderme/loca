package main

import (
	"regexp"
)

var (
	reHREF      = regexp.MustCompile(`(?i) href=["']([^<>"']+)['"]`)
	reSRC       = regexp.MustCompile(`(?i) src=["']([^<>"']+)['"]`)
	reUglyURL   = regexp.MustCompile(`(?s)[;\+\s"'` + "`" + `%\?!~><\}\{\]\[\\\\:,|\*&^$@\+\(\)]+`)
	reOneHyphen = regexp.MustCompile(`-{2,}`)
	// integrety and crossorigin maybe problematic on localhost
	reCleanAttrs = []*regexp.Regexp{
		// crossorigin="anonymous"
		regexp.MustCompile(`(?is)(<[^<>]+)(\scrossorigin=["'][^"'<>]*["'])`),
		// sintegrity="sha256-nas+AAA/ac129be5a99129c50ec2bf891b2b439ac02bedde2a73c6e91f683f79c43577b6"
		regexp.MustCompile(`(?is)(<[^<>]+)(\sintegrity=["'][^"'<>]*["'])`),
		regexp.MustCompile(`(?is)(<[^<>]+)(\stype=["']\s*text/javascript\s*["'])`),
	}
	reWhitespace = regexp.MustCompile(`(?s)\s+`)
)
