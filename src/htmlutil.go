package main

import (
	"bytes"
	html5 "golang.org/x/net/html"
	//	"html"
	"strings"
)

func cleanHTML(s string) (string, error) {
	data := strings.NewReader(s)
	buff := &bytes.Buffer{}

	node, err := html5.Parse(data)

	if err != nil {
		return "", err
	}

	err = html5.Render(buff, node)

	return buff.String(), err

}

func packHTML(s string) string {
	s = reWhitespace.ReplaceAllLiteralString(s, " ")
	s = strings.TrimSpace(s)

	return s
}

func cleanAttrs(s string) string {
	for _, r := range reCleanAttrs {
		s = r.ReplaceAllString(s, "$1")
	}

	return s
}

func prettyURL(ugly string) (pretty string) {
	ugly = reUglyURL.ReplaceAllLiteralString(ugly, "-")
	pretty = reOneHyphen.ReplaceAllLiteralString(ugly, "-")

	// what about basename?
	return

}
