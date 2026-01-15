package main

import (
	"testing"

	"github.com/ross-spencer/zenodocfl/internal/types"
)

type reltest struct {
	filename string
	base     base
}

var recordTests = []reltest{
	{
		filename: "testdata/c02.json",
		base: base{
			Title: []title{
				{"en", "C02 Beata progenies"},
			},
			Signature: "",
			License:   "",
			Publisher: "",
			Poster:    poster{},
		},
	},
	{
		filename: "testdata/m001.json",
		base: base{
			Title: []title{
				{"en", "M001 Beata progenies"},
			},
			Signature: "",
			License:   "",
			Publisher: "",
			Poster:    poster{},
		},
	},
}

func TestRecords(t *testing.T) {
	for _, test := range recordTests {
		record, _ := readJSON(test.filename)
		if record.Base.Title[0].Value != test.base.Title[0].Value {
			t.Errorf("title values are incorrect")
		}
		// NB. more tests can be added the more involved the data gets.
	}
}

func TestConvertURI(t *testing.T) {
	input := "mediaserver:hsm/motet_cycles_data_motet_cycles_data_MEI_files_motets_M001BeataProgenies.png"
	output := "https://ba14ns21403-sec1.fhnw.ch/mediasrv/hsm/motet_cycles_data_motet_cycles_data_MEI_files_motets_M001BeataProgenies.png/master"
	res := convertMediaServerURI(input)
	if res != output {
		t.Errorf("uri not converted correctly: '%s' expected: '%s'", res, output)
	}
}

func TestConvertIdentifier(t *testing.T) {
	inputKey := "ark"
	inputValue := "ark:/15737/p655-cz83-8mxn"
	outputValue := "https://n2t.net/ark:/15737/p655-cz83-8mxn"
	res := convertIdentifier(inputKey, inputValue)
	if res != outputValue {
		t.Errorf("ark conversion failed: '%s' expected: '%s'", res, outputValue)
	}
	inputKey = "handle"
	inputValue = "20.500.11806/med/3h1x-65hq-1w"
	outputValue = "https://hdl.handle.net/20.500.11806/med/3h1x-65hq-1w"
	res = convertIdentifier(inputKey, inputValue)
	if res != outputValue {
		t.Errorf("handle conversion failed: '%s' expected: '%s'", res, outputValue)
	}
}

func makeTestItem() types.Item {
	item := types.Item{}

	item.Poster.Url = "http://example.com/3"

	item.Relationship = []types.Relationship{
		{
			Label:  "",
			Url:    "http://example.com/1",
			Poster: types.Poster{},
		},
	}

	item.Relationship[0].Poster.Url = "http://example.com/4"

	item.Media = []types.Media{
		{
			Name:     "",
			MimeType: "",
			Url:      "http://example.com/2",
		},
	}

	return item
}

func TestGetURLS(t *testing.T) {
	collection := types.Collection{}
	items := []types.Item{
		makeTestItem(),
		makeTestItem(),
	}
	collection.Items = items
	collection.GetURLs()
	if len(collection.PosterURLs) != 2 {
		t.Errorf("urls length should be two: %d", len(collection.PosterURLs))
	}
}
