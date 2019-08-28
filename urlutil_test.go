package main

import (
	"testing"
)

type URLFormat struct {
	Parent,
	Child,
	Expected string
	Unescape bool
}

type URLRelation struct {
	Parent,
	Child string
	Result bool
}

func TestParseURL(t *testing.T) {
	URLs := map[string]bool{
		"example.com/home.html":                   true,
		"https://example.com/list.htm?v=1":        true,
		"http://example.com/item.htm?v=1&lang=en": true,

		"ftps://example.com/list.htm?v=1":        false,
		"ftp://example.com/item.htm?v=1&lang=en": false,
	}

	for u, valid := range URLs {
		_, err := parseURL(u)
		if valid && err != nil {
			t.Error(u, "is valid, but got err", err)
		}

		if !valid && err == nil {
			t.Error(u, "is invalid, but got no err")
		}

	}

}

func TestDiscoverSrcURLs(t *testing.T) {
	const (
		s = `
<img src="https://example.com/imgs/logo.png" alt="Logo">
<img src="" width=1 height=2>
<SCRIPT SRC='https://example.com/assests/js/lib.js'></SCRIPT>

`
		logo   = "https://example.com/imgs/logo.png"
		script = "https://example.com/assests/js/lib.js"
	)
	URLs := discoverSrcURLs(s)

	if len(URLs) != 2 {
		t.Error("Wrong discovered URLs count")
		return
	}

	if URLs[0] != logo {
		t.Error(URLs[0], "doesn't equal", logo)
	}

	if URLs[1] != script {
		t.Error(URLs[0], "doesn't equal", script)
	}

}

func TestDiscoverHREFURLs(t *testing.T) {
	const (
		s = `

<a href="https://example.com/home"></a>
<a href="https://example.com/home?lang=en"></a>
<a HREF='https://example.com/faqs.html'></A>

`
		home   = "https://example.com/home"
		homeEn = "https://example.com/home?lang=en"
		faqs   = "https://example.com/faqs.html"
	)
	URLs := discoverHREFURLs(s)

	if len(URLs) != 3 {
		t.Error("Wrong discovered URLs count")
		return
	}

	if URLs[0] != home {
		t.Error(URLs[0], "doesn't equal", home)
	}

	if URLs[1] != homeEn {
		t.Error(URLs[0], "doesn't equal", homeEn)
	}

	if URLs[2] != faqs {
		t.Error(URLs[0], "doesn't equal", faqs)
	}

}

func TestResolveURL(t *testing.T) {

	URLs := []URLFormat{
		URLFormat{
			Parent:   "https://www.example.com",
			Child:    "home.html",
			Expected: "https://www.example.com/home.html",
			Unescape: false,
		},

		URLFormat{
			Parent:   "https://www.example.com/",
			Child:    "../../home.html?lang=ja",
			Expected: "https://www.example.com/home.html?lang=ja",
			Unescape: false,
		},
	}

	for _, u := range URLs {
		resolved, _ := resolveURL(u.Parent, u.Child, u.Unescape)

		if u.Expected != resolved {
			t.Error("Expected:", u.Expected, "But Got:", resolved)
		}
	}

}

func TestResolvePath(t *testing.T) {

	paths := []URLFormat{
		URLFormat{
			Parent:   "/cat/list.html",
			Child:    "../../images/logo.png",
			Expected: "/images/logo.png",
		},

		URLFormat{
			Parent:   "/cat/subcats/list.html",
			Child:    "../images/logo.png",
			Expected: "/cat/images/logo.png",
		},

		URLFormat{
			Parent:   "/cat/subcats/list.html",
			Child:    "./images/logo.png",
			Expected: "/cat/subcats/images/logo.png",
		},

		URLFormat{
			Parent:   "/cat/subcats/list.html",
			Child:    "images/logo.png",
			Expected: "/cat/subcats/images/logo.png",
		},
	}

	for _, p := range paths {
		resolved := resolvePath(p.Parent, p.Child)
		if p.Expected != resolved {
			t.Error("Expected:", p.Expected, "But Got:", resolved)
		}
	}

}

func TestIsParent(t *testing.T) {

	urls := []URLRelation{

		URLRelation{
			Parent: "https://example.org/about",
			Child:  "https://example.org/about/site",
			Result: true,
		},

		URLRelation{
			Parent: "https://example.org/home",
			Child:  "https://example.org/about/",
			Result: false,
		},

		URLRelation{
			Parent: "https://example.org/home/",
			Child:  "https://example.org/about/",
			Result: false,
		},

		URLRelation{
			Parent: "https://example.org/home.html",
			Child:  "https://example.org/about/",
			Result: false,
		},
	}

	for _, u := range urls {

		result := isParent(u.Parent, u.Child)

		if u.Result != result {

			t.Error("Expected:", u.Result, "But Got:", result)

		}

	}

}
