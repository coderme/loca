package main

import (
	"regexp"
)

var (
	reHREF     = regexp.MustCompile(`(?i) href=["']([^<>"']+)['"]`)
	reStyle    = regexp.MustCompile(`(?i)<link\s+rel=['"]stylesheet['"]\s+href=['"]([^<>"']+)['"]`)
	reStyleAlt = regexp.MustCompile(`(?i)<link\s+href=['"]([^<>"']+)['"]\s+rel=['"]stylesheet['"]`)
	reSRC      = regexp.MustCompile(`(?i) src=["']([^<>"']+)['"]`)

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

	reCleanTags = []*regexp.Regexp{

		regexp.MustCompile(`(?is)<meta[^<>]+rel="pingback"[^<>]*>`),
	}

	reFilter = []*regexp.Regexp{
		regexp.MustCompile(`(?is)(src=")([^<>"]+)(")`),
		regexp.MustCompile(`(?is)(src=')([^<>']+)(')`),
		regexp.MustCompile(`(?is)(href=")([^<>"]+)(")`),
		regexp.MustCompile(`(?is)(href=')([^<>"]+)(')`),
	}

	reWhitespace = regexp.MustCompile(`(?s)\s+`)
)
