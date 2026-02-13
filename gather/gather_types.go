// types describe the data structures we anticipate in the JSON data
// retrieve from the INK server.

package main

type title struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

type base struct {
	// Title of the record.
	Title []title `json:"title"`
	// Signature / slug of the record.
	Signature string `json:"signature"`
	// License belonging to the item.
	License string `json:"license"`
	// Publisher is the item's publisher.
	Publisher string `json:"publisher"`
	// Poster belonging to the main item.
	Poster poster `json:"poster"`
}

type poster struct {
	Name string `json:"name" `
	Url  string `json:"uri"`
}

type references struct {
	Signature string  `json:"signature"`
	Title     []title `json:"title"`
	Url       string  `json:"url"`
	Poster    poster  `json:"poster"`
	License   string  `json:"license"`
	Media     media   `json:"media"`
}

type item struct {
	Name string `json:"name"`
	Url  string `json:"uri"`
	Mime string `json:"mimetype"`
}

type media struct {
	Mediatype string `json:"type"`
	Items     []item `json:"items"`
}

type extra struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type note struct {
	Text  string `json:"text"`
	Title string `json:"title"`
}

type inkRecord struct {
	// Core metadata.
	Base base `json:"base"`
	// The fileName queried.
	FileName string `json:"file_name"`
	// Media associated with the record.
	Media []media `json:"media"`
	// Relationships to the item.
	ReferencesFull []references `json:"referencesFull"`
	// Extra metadata referenced in the item, usually
	// persistent identifiers.
	Extra []extra `json:"extra"`
	// Notes contains information like Description.
	Notes []note `json:"notes"`
	// Source data used to create this record.
	Source string `json:"source"`
}
