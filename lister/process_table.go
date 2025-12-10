package main

import (
	"fmt"
	"strings"

	"github.com/ross-spencer/zenodocfl/internal/types"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func makeDataURL(url string) string {
	url = strings.Replace(url, "/detail/", "/detailjson/", 1)
	// A quirk of the INK system is that it must end with the lang
	// param. This can probably be ironed out upstream.
	return fmt.Sprintf("%s%s", url, "de?")
}

// getImg provides a helper function for normalizing preview images
// returned in Katalog results.
func getImg(mediaLink string) string {
	const previewIMGURLSlug string = "/resize/size100x100/formatPNG/autorotate"
	return strings.TrimSpace(strings.Split(mediaLink, previewIMGURLSlug)[0])
}

// processTableData will, given a html.Node table “tBody` process the
// and extract mediaThekRecord results.
func processTableData(n *html.Node, m *[]types.MediathekRecord, lang string) {

	const tdImg int = 0
	const tdSig int = 1
	const tdTitel int = 2
	const tdYear int = 3
	const tdArtist int = 4
	const tdEvent int = 5

	if n.Type == html.ElementNode && n.DataAtom == atom.Tr {
		r := types.MediathekRecord{}
		for _, a := range n.Attr {
			if a.Key == "onclick" {
				s := a.Val
				s = strings.Replace(s, "window.location='", "", 1)
				s = strings.Replace(s, "\\/\\/", "//", 1)
				s = strings.Split(s, fmt.Sprintf("%s?", lang))[0]
				r.Url = s
				r.DataURL = makeDataURL(s)
				break
			}
		}
		idx := 0
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			for f := c.FirstChild; f != nil; f = f.NextSibling {
				switch val := idx; val {
				case tdImg:
					for _, s := range f.Attr {
						if s.Key == "src" {
							r.Img = getImg(s.Val)
							break
						}
					}
				case tdTitel:
					r.Title = strings.TrimSpace(f.Data)
				case tdSig:
					r.Signature = strings.TrimSpace(f.Data)
				default:
					// not yet handled...
				}
			}
			if c.DataAtom != 0 {
				idx += 1
			}
		}
		*m = append(*m, r)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		processTableData(c, m, lang)
	}
}

// processTableResults recursively process html.Node data to extract
// table data from Katalog results and forwards that for processing.
func processTableResults(n *html.Node, m *[]types.MediathekRecord, lang string) {
	if n.Type == html.ElementNode && n.DataAtom == atom.Tbody {
		// TODO... all processing happens here...
		processTableData(n, m, lang)
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		processTableResults(c, m, lang)
	}
}
