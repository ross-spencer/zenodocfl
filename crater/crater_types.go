package main

type rocrate struct {
	Context string        `json:"@context"`
	Graph   []interface{} `json:"@graph"`
}

type idPointer struct {
	ID string `json:"@id,omitempty"`
}

type root struct {
	ID         string    `json:"@id"`
	Type       string    `json:"@type,omitempty"`
	ConformsTo idPointer `json:"conformsTo,omitempty"`
	Identifier string    `json:"identifier,omitempty"`
	About      idPointer `json:"about,omitempty"`
}

type files struct {
	ID            string      `json:"@id"`
	Name          string      `json:"name,omitempty"`
	Type          string      `json:"@type,omitempty"`
	ContentURL    string      `json:"contentUrl,omitempty"`
	DatePublished string      `json:"datePublished,omitempty"`
	Description   string      `json:"description,omitempty"`
	HasPart       []idPointer `json:"hasPart,omitempty"`
	Identifier    string      `json:"identifier,omitempty"`
	Keywords      []string    `json:"keywords,omitempty"`
	License       string      `json:"license,omitempty"`
	Publisher     idPointer   `json:"publisher,omitempty"`
}

type org struct {
	ID   string `json:"@id"`
	Type string `json:"@type,omitempty"`
	Name string `json:"name,omitempty"`
}
