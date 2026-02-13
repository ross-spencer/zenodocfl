package main

type rocrate struct {
	Context string  `json:"@context"`
	Graph   []graph `json:"@graph"`
}

type idPointer struct {
	ID string `json:"@id,omitempty"`
}

type graph struct {
	ID            string      `json:"@id"`
	Name          string      `json:"name,omitempty"`
	Type          string      `json:"@type,omitempty"`
	About         idPointer   `json:"about,omitempty"`
	ContentURL    string      `json:"contentUrl,omitempty"`
	DatePublished string      `json:"datePublished,omitempty"`
	Description   string      `json:"description,omitempty"`
	HasPart       []idPointer `json:"hasPart,omitempty"`
	Identifier    string      `json:"identifier,omitempty"`
	Keywords      []string    `json:"keywords,omitempty"`
	License       string      `json:"license,omitempty"`
	Publisher     idPointer   `json:"publisher,omitempty"`
}
