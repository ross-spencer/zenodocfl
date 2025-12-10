/* types used across the INK -> Zenodo workflow. */

package types

import (
	"fmt"
	"log"
	"slices"
)

/* Lister types */

// mediathekRecord will be part of our builder pattern that enables
// us to build a complete record we can then use to create a
// Zenodo dataset.
type MediathekRecord struct {
	// URL of the INK record.
	Url string `json:"url"`
	// Slug belonging to the record.
	Signature string `json:"signature"`
	// Image associated with the record.
	Img string `json:"img"`
	// Title of the record.
	Title string `json:"title"`
	// Data record associated with record, usually JSON.
	DataURL string `json:"dataUrl"`
}

// String(er) for the mediathekRecord.
func (r MediathekRecord) String() string {
	return fmt.Sprintf("title: %s: | url: %s", r.Title, r.Url)
}

/* Gather types */

// Item serves to flatten the INK record into something that starts
// to look more like the RO-CRATE record we will create.
type Item struct {
	// The item title.
	Label string `json:"label"`
	// The file used to create the record.
	File string `json:"file"`
	// License belonging to the item.
	License string `json:"license"`
	// Publisher of the item.
	Publisher string `json:"publisher"`
	// Relationship describes relationships to an item one way or
	// another. The direction of relationship is not yet
	// determined.
	//
	// NB. likely hasPart in RO-CRATE.
	Relationship []Relationship `json:"relationships"`
	// Meedia describes media associated with a record.
	//
	// NB. likely hasFile in RO-CRATE.
	Media []Media `json:"media,omitempty"`
	// Identifers associated with the record.
	Identifiers []Identifier `json:"identifiers"`
	// Poster provides an image associated with a record. The basee
	// record has a poster and relations often do as well.
	Poster Poster `json:"poster"`
	// Description describes the record.
	Description string `json:"description,omitempty"`
}

type Relationship struct {
	Label  string `json:"label"`
	Url    string `json:"url"`
	Poster Poster `json:"poster"`
}

type Media struct {
	Name     string `json:"name"`
	MimeType string `json:"mimetype"`
	Url      string `json:"url"`
}

type Identifier struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Poster struct {
	Name string
	Url  string `json:"url"`
}

// Collection provides a type as a mechanism for bringing record and
// relationship information together.
/* Example:

   "collection": [
        // record1
        // record2
   ],
   "all_uurls": [
       "url1",
       "url2"
   ]

*/
type Collection struct {
	Items      []Item `json:"records"`
	itemURLs   []string
	MediaURLs  []string `json:"media_urls"`
	PosterURLs []string `json:"poster_urls"`
}

// addItem determines if an item should be added to a slice for
// the purpose of creating a manifest.
func addItem(list []string, item string) bool {
	if item == "" {
		return false
	}
	if slices.Contains(list, item) {
		return false
	}
	return true
}

// Count the number of URLs in the collection.
func (collection *Collection) GetURLs() {
	itemUrls := []string{}
	mediaUrls := []string{}
	posterUrls := []string{}
	for _, item := range collection.Items {
		if addItem(posterUrls, item.Poster.Url) {
			posterUrls = append(posterUrls, item.Poster.Url)
		}
		for _, rel := range item.Relationship {
			if addItem(posterUrls, rel.Poster.Url) {
				posterUrls = append(posterUrls, rel.Poster.Url)
			}
			if !addItem(itemUrls, rel.Url) {
				continue
			}
			itemUrls = append(itemUrls, rel.Url)
		}
		for _, med := range item.Media {
			if !addItem(mediaUrls, med.Url) {
				continue
			}
			mediaUrls = append(mediaUrls, med.Url)
		}
	}
	collection.itemURLs = itemUrls
	collection.MediaURLs = mediaUrls
	collection.PosterURLs = posterUrls

	log.Printf(
		"itemURLs: '%d', medialURLs: '%d', posterURLs: '%d'",
		len(collection.itemURLs),
		len(collection.MediaURLs),
		len(collection.PosterURLs),
	)
}
