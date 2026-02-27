package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ross-spencer/zenodocfl/internal/types"
)

func readManifest(manifest string) types.Collection {
	data, err := os.ReadFile(manifest)
	if err != nil {
		log.Println("error reading collection manifest:", err)
		os.Exit(1)
	}
	var collection types.Collection
	err = json.Unmarshal(data, &collection)
	if err != nil {
		log.Println("error reading collection manifest:", err)
		os.Exit(1)
	}
	return collection

}

// downloadCrateObj enables us to download material from a given URL
// and save ti in the given folder.
func downloadCrateObj(url string, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating path: %w (%s)", err, path)
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading url: %w (%s)", err, url)
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error accessing url data: %w (%s)", err, url)
	}
	return nil
}

// createCrateObj handles the creation of a RO-Crate object based on
// the information provided by the user.
func createCrateObj(path string, data string) {
	log.Println("metadata path:", path)
	log.Printf("preview with: `rochtml %s` (if installed)", path)
	err := os.WriteFile(path, []byte(fmt.Sprintf("%s\n", data)), 0755)
	if err != nil {
		log.Println("unable to write to file;", err)
	}
}

// moveRecords addds the collection records to the RO-CRATE.
func moveRecords(records []types.Item, path string, partPrefix string) []string {
	parts := []string{}
	for _, item := range records {
		filePath := filepath.Join(path, item.File)
		err := os.WriteFile(filePath, []byte(fmt.Sprintf("%s\n", item.Source)), 0755)
		if err != nil {
			log.Println("unable to write to file;", err)
		}
		parts = append(parts, fmt.Sprintf("%s/%s", partPrefix, item.File))
	}
	return parts
}

// makeFIlename returns a filename for the URL we're downloading.
func makeFilename(url string) string {
	url = strings.Replace(url, "$$poster/master", "", 1)
	url = strings.Replace(url, "/master", "", 1)
	split := strings.Split(url, "/")
	return split[len(split)-1]
}

// downloadFile retrieves media from the server and stores it in the
// given path.
func downloadFile(urls []string, path string, partPrefix string, dryrun bool) []string {
	parts := []string{}
	for _, url := range urls {
		fileName := makeFilename(url)
		filePath := filepath.Join(path, fileName)
		if debug {
			log.Println(filePath)
		}

		if !dryrun {
			err := downloadCrateObj(url, filePath)
			if err != nil {
				log.Printf("cannot download object: %s", err)
				os.Exit(1)
			}
		}
		parts = append(parts, fmt.Sprintf("%s/%s", partPrefix, fileName))
	}
	return parts
}

const crateName string = "ro-crate-metadata.json"
const creativeWork string = "CreativeWork"
const rocrateContext string = "https://w3id.org/ro/crate/1.1/context"
const rocrateConform string = "https://w3id.org/ro/crate/1.1"

type userData struct {
	Identifier    string `json:"identifier"`
	Description   string `json:"description"`
	Name          string `json:"name"`
	RecordType    string `json:"type"`
	DatePublished string `json:"data_published"`
	License       string `json:"license"`
	Keywords      string `json:"keywords"`
	Publisher     string `json:"publisher"`
	PublisherName string `json:"publisherName"`
	// we might not always have a canonical url.
	Url string `json:"url"`
	// added automatically.
	parts []string
}

func (userData userData) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
		userData.Identifier,
		userData.Description,
		userData.Name,
		userData.RecordType,
		userData.DatePublished,
		userData.License,
		userData.Keywords,
		userData.Publisher,
		userData.PublisherName,
		userData.Url,
	)
}

func getKeywords(values string) []string {
	keywords := []string{}
	tmpKeys := strings.Split(values, ",")
	for _, item := range tmpKeys {
		keywords = append(keywords, strings.TrimSpace(item))
	}
	return keywords
}

// makeCrateObj creates a RO-CRATE JSON object.
func makeCrateObj(userData userData) rocrate {

	const rootID string = "./"
	const orgType string = "Organization"

	crate := rocrate{}
	crate.Context = rocrateContext
	meta := root{}
	meta.ID = crateName
	meta.Identifier = crateName
	meta.Type = creativeWork
	meta.About = idPointer{rootID}
	meta.ConformsTo = idPointer{rocrateConform}
	obj := files{}
	obj.ID = rootID
	obj.Identifier = userData.Identifier
	obj.Type = userData.RecordType
	obj.Name = userData.Name
	obj.Description = userData.Description
	obj.License = userData.License
	obj.DatePublished = userData.DatePublished
	obj.Publisher = idPointer{userData.Publisher}
	obj.Keywords = getKeywords(userData.Keywords)
	obj.ContentURL = userData.Url
	for _, item := range userData.parts {
		obj.HasPart = append(obj.HasPart, idPointer{item})
	}
	crate.Graph = append(crate.Graph, meta)
	crate.Graph = append(crate.Graph, obj)
	pub := org{}
	pub.Name = userData.PublisherName
	pub.ID = userData.Publisher
	pub.Type = orgType
	crate.Graph = append(crate.Graph, pub)
	return crate
}
