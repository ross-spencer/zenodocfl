package main

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/ross-spencer/zenodocfl/internal/types"
)

type testData struct {
	name         string
	path         string
	lang         string
	len          int
	sampleIndex  int
	sampleResult types.MediathekRecord
}

var extractTests = []testData{
	{
		"motet-table",
		"testdata/results.htm",
		"en",
		271,
		50,
		types.MediathekRecord{
			Url:     "https://ink.sammlung.cc/detail/motetcycle-0445/",
			Title:   "C55 Gaude flore virginali",
			DataURL: "https://ink.sammlung.cc/detailjson/motetcycle-0445/de?",
		},
	},
	{
		"motet-summary",
		"testdata/summary_results.htm",
		"en",
		59,
		11,
		types.MediathekRecord{
			Url:     "https://ink.sammlung.cc/detail/motetcycle-0409/",
			Title:   "C12a Ave mundi domina",
			DataURL: "https://ink.sammlung.cc/detailjson/motetcycle-0409/de?",
		},
	},
	{
		"motet-summary",
		"testdata/summary_de.htm",
		"de",
		59,
		11,
		types.MediathekRecord{
			Url:     "https://ink.sammlung.cc/detail/motetcycle-0409/",
			Title:   "C12a Ave mundi domina",
			DataURL: "https://ink.sammlung.cc/detailjson/motetcycle-0409/de?",
		},
	},
}

// TestTableExtract ensures that the code we use to extract data from
// the Unibas Katalog is extracted correctly.
func TestTableExtract(t *testing.T) {

	var err error
	m := []types.MediathekRecord{}
	for _, test := range extractTests {
		testData, _ := os.Open(test.path)
		defer testData.Close()
		reader := bufio.NewReader(testData)
		m, err = extractTable(reader, test.lang)
		if err != nil {
			t.Errorf("unexpected error in test: '%s' (%s)", test.name, err)
		}
		if len(m) != test.len {
			t.Errorf("incorrect number of results from test: '%s', %d expected %d", test.name, len(m), test.len)
		}
		compare := m[test.sampleIndex]
		if compare.Url != test.sampleResult.Url {

			t.Errorf("%s", compare.Url)
		}
		if compare.Title != test.sampleResult.Title {
			t.Errorf("%s", compare.Title)
		}
		if compare.DataURL != test.sampleResult.DataURL {
			t.Errorf("%s", compare.DataURL)
		}
		if !strings.HasSuffix(compare.DataURL, lang) {
			t.Errorf("data URL must end with language param (%s): %s", lang, compare.DataURL)
		}
	}
}

// TestMakeINKURL ensures that the INK URL is created accurately.
func TestMakeINKURL(t *testing.T) {
	expected := "https://ink.sammlung.cc/table/fr?search=motetcycle&cursor=eyJmcm9tIjogMCwgInNpemUiOiA1fQ=="
	url := makeINKURL("fr", "motetcycle", 5)
	if url != expected {
		t.Errorf("url generation failing: '%s', expected: '%s'", url, expected)
	}
}
