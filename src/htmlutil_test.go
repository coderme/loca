package main

import (
	"strings"
	"testing"
)

func TestPackHTML(t *testing.T) {
	const html = `   

<!DOCTYPE html>
          <head>
<title>   Sample   </title>

</head>
      <body>


    Content&nbsp;&nbsp;&nbsp;&nbsp;

</body>

  </html>

`

	packed := packHTML(html)

	checks := map[string]string{
		"  ":   "double whitespace",
		"\n\n": "double newlines",
		" \n":  "space followed by a newline",
		"\t ":  "tab followed by a space",
		" \t":  "space followed by atab",
	}

	for s, m := range checks {

		if strings.Contains(packed, s) {
			t.Error("packHTML() result contains " + m)
		}

	}

}
