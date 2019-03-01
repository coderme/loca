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
			t.Error("packHTML() result contains", m)
		}

	}

}

func TestCleanAttrs(t *testing.T) {
	const snippet = `
<link href="https://www.example.com/pretty.css?refresh=true" rel="stylesheet">
  <link rel="stylesheet" type="text/css" integrity="sha256-TN63maDMYkAimHYRvtXvanX7sCK5Sf7ZCoZA/YAlKFQ=" crossorigin="anonymous" href="/assets/main.css">
  
  <script src="https://cdn.example.com/lib.min.js"></script>
  <script src="/assets/main.js" integrity="sha256-nas+ETXQ27ZXx/j+QlYE5VFONJwqbqoBVLVAHGuKIiE=" crossorigin="anonymous" type="text/javascript"></script>

`

	checks := map[string]string{
		`crossorigin=`: "crossorigin attribute",
		`anonymous=`:   "anonymous attribute",
	}

	cleaned := strings.ToLower(
		cleanAttrs(snippet),
	)

	for s, m := range checks {
		if strings.Contains(cleaned, s) {
			t.Error(`cleanAttrs() contains`, m)
		}

	}

}

func TestPrettyURL(t *testing.T) {
	const (
		uglyChars = `"'<> ;[]{}%~,&|*$@+()!` + "`"
		ugly      = "https://www.example.com/list.html?page=13&order=desc&note=" + uglyChars
	)

	pretty := prettyURL(ugly)

	if strings.ContainsAny(pretty, uglyChars) {
		t.Error(pretty, "contains at least on chars of",
			uglyChars,
		)
	}

}
