/*
gather is responsible for bulk downloading urls, relationships and
auxiliary objects. The output be used to build a RO-CRATE gocfl package.

 1. access our current manifest describing the INK records.
 2. Read all the data.
 3. Using the data create a flattened version to be fed into `crater`.
 4. Grab all URLs that are likely to make up the final record.
 5. Output a summary report for crater.
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ross-spencer/zenodocfl/internal/logformatter"
	"github.com/ross-spencer/zenodocfl/internal/types"
)

var (
	download string
	list     bool
	vers     bool

	// app constants.
	version = "dev-0.0.0"
	commit  = "000000000000000000000000000000000baddeed"
	date    = "1970-01-01T00:00:01Z"
)

var agent string = fmt.Sprintf("INK-gather/%s", version)

// initFlags initializes the flags we use with this app.
func initFlags() {
	flag.StringVar(&download, "download", "", "download items in the given manifest")
	flag.BoolVar(&list, "list", false, "list records in the JSON directoru already downloaded")
	flag.BoolVar(&vers, "version", false, "Return version")
}

// prettyJSON outputs prettified JSON.
func prettyJSON(content []byte) ([]byte, error) {
	var pretty interface{}
	err := json.Unmarshal(content, &pretty)
	if err != nil {
		return []byte{}, err
	}
	prettyJSON, err := json.MarshalIndent(pretty, "", " ")
	if err != nil {
		return []byte{}, err
	}
	return prettyJSON, nil
}

// downloadManifest will download the data files associated with the
// input manifest.
func downloadManifest(dl string) []types.MediathekRecord {
	log.Println("downloading from:", dl)
	data, err := os.ReadFile(dl)
	if err != nil {
		log.Println("error reading lister manifest:", err)
		os.Exit(1)
	}
	splitPaths := strings.Split(string(data), "\n")
	paths := []types.MediathekRecord{}
	for _, v := range splitPaths {
		if v == "" {
			continue
		}
		var record types.MediathekRecord
		json.Unmarshal([]byte(v), &record)
		paths = append(paths, record)
	}
	log.Println("records to download:", len(paths))
	return paths
}

// downloadFiles downloads the files from the manifest into the
// data folder.
func downloadFiles(files []types.MediathekRecord) {
	for _, v := range files {
		log.Println("downloading;", v.DataURL)
		resp, err := http.Get(v.DataURL)
		if err != nil {
			log.Printf("network error reading data file: '%s' (%s)", v.DataURL, err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Println("status code != 200:", resp.StatusCode)
			continue
		}
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error reading data file: '%s' (%s)", v.DataURL, err)
			continue
		}
		prettified, _ := prettyJSON(content)
		path := filepath.Join("data", fmt.Sprintf("%s.json", v.Signature))
		os.Remove(path)
		err = os.WriteFile(path, prettified, 0644)
		if err != nil {
			log.Println("unable to write to file;", v.Signature)
		}
		time.Sleep(1 * time.Second)
	}
}

func getTitle(title []title) (string, error) {
	if len(title) < 1 {
		return "", fmt.Errorf("title is empty")
	}
	return title[0].Value, nil
}

func getDescription(notes []note) (string, error) {
	if len(notes) < 1 {
		return "", fmt.Errorf("no notes associated with record")
	}
	description := ""
	for _, v := range notes {
		if v.Title != "Description" {
			continue
		}
		description = v.Text
		break
	}
	description = strings.Replace(description, "<p>", "", 1)
	description = strings.Replace(description, "</p>", "", 1)
	return description, nil
}

// addItemMD creates the primary record metadata to translate to
// top-level items in RO-CRATE.
func addItemMD(record inkRecord) (types.Item, error) {
	item := types.Item{}
	title, err := getTitle(record.Base.Title)
	if err != nil {
		log.Println("cannot retrieve title from record")
		return item, fmt.Errorf("cannot retrieve title from record")
	}
	item.Label = title
	item.File = record.FileName
	item.License = record.Base.License
	item.Publisher = record.Base.Publisher
	item.Poster.Name = record.Base.Poster.Name
	item.Poster.Url = convertMediaServerURI(record.Base.Poster.Url)
	description, err := getDescription(record.Notes)
	if err != nil {
		log.Println("cannot retrieve description from record")
	}
	item.Description = description
	return item, nil
}

// addItemRelationships associates relationships with the new
// record. This is likely hasPart in RO-CRATE - TBD.
func addItemRelationships(item types.Item, record inkRecord) (types.Item, error) {
	rels := []types.Relationship{}
	for _, value := range record.ReferencesFull {
		rel := types.Relationship{}
		rel.Poster.Name = value.Poster.Name
		rel.Poster.Url = convertMediaServerURI(value.Poster.Url)
		title, err := getTitle(value.Title)
		if err != nil {
			log.Println("cannot retrieve title from record")
			return item, fmt.Errorf("cannot retrieve relationship title from record")
		}
		rel.Label = title
		rel.Url = value.Url
		rels = append(rels, rel)
	}
	item.Relationship = rels
	return item, nil
}

// addMedia associates media objects with the record. This is likely
// hasFle in RO-CRATE - TBD.
func addMedia(item types.Item, record inkRecord) types.Item {
	meds := []types.Media{}
	for _, value := range record.Media {
		for _, mediaItem := range value.Items {
			med := types.Media{}
			med.Name = mediaItem.Name
			med.MimeType = mediaItem.Mime
			med.Url = convertMediaServerURI(mediaItem.Url)
			meds = append(meds, med)
		}
	}
	item.Media = meds
	return item
}

// addIdentifiers ensures identifiers are associated with the
// record. These will all be displayed in the RO-CRATE.
func addIdentifiers(item types.Item, record inkRecord) types.Item {
	ids := []types.Identifier{}
	for _, value := range record.Extra {
		id := types.Identifier{}
		switch value.Key {
		case "ark":
		case "handle":
			id.Type = value.Key
			id.Name = value.Value
			id.Url = convertIdentifier(value.Key, value.Value)
			ids = append(ids, id)
		default:
			continue
		}
	}
	item.Identifiers = ids
	return item
}

// makeCollection returns a more complete collection manifest that can
// be given to crater to create a RO-CRATE package.
func makeCollection(manifest []inkRecord) types.Collection {
	collection := types.Collection{}
	for _, record := range manifest {
		item, err := addItemMD(record)
		if err != nil {
			log.Println("cannot retrieve title for item")
			continue
		}
		item, err = addItemRelationships(item, record)
		item = addMedia(item, record)
		item = addIdentifiers(item, record)
		collection.Items = append(collection.Items, item)
	}

	// gather all remaining URLs for download.
	collection.GetURLs()
	return collection
}

// listJSON will output a slice of all the records associated with
// the given manifest. The data directory should already exist to
// enable this.
func listJSON() []inkRecord {
	const dataDir string = "data"
	const dataExt string = "json"
	manifest := []inkRecord{}
	if !exists(dataDir) {
		log.Println("'data' directory doesn't exist")
		os.Exit(1)
	}
	entries, err := os.ReadDir(dataDir)
	log.Println("items in data directory:", len(entries))
	if err != nil {
		log.Println("an unknown error has occurred:", err)
	}
	for _, entry := range entries {
		fname := entry.Name()
		if strings.HasPrefix(fname, ".") {
			// file is hidden and shouldn't be processed.
			continue
		}
		if !strings.HasSuffix(fname, dataExt) {
			// file isn't a valid data file.
			continue
		}
		filePath := fmt.Sprintf("%s%c%s", dataDir, os.PathSeparator, fname)
		log.Println("processing;", filePath)
		record, err := readJSON(filePath)
		if err != nil {
			log.Println("error processing data:", err)
			continue
		}
		record.FileName = fname
		manifest = append(manifest, record)
	}
	return manifest
}

func main() {

	logformatter.Set("gather", true)

	initFlags()
	flag.Parse()

	if vers {
		fmt.Fprintf(os.Stderr, "%s (%s) commit: %s date: %s\n", agent, version, commit, date)
		os.Exit(0)
	} else if flag.NFlag() < 1 {
		fmt.Fprintln(os.Stderr, "Usage:  ")
		fmt.Fprintln(os.Stderr, "        OPTIONAL: [-download]  STRING")
		fmt.Fprintln(os.Stderr, "        OPTIONAL: [-list] ")
		fmt.Fprintln(os.Stderr, "        OPTIONAL: [-version] ")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Output: [STRING] {result JSON}")
		fmt.Fprintf(os.Stderr, "Output: [STRING] {version: '%s'}\n\n", agent)
		flag.Usage()
		os.Exit(0)
	}

	if vers {
		fmt.Println(agent)
		return
	}

	if download != "" {
		files := downloadManifest(download)
		downloadFiles(files)
		return
	}

	if list {
		manifest := listJSON()
		collection := makeCollection(manifest)
		printCollection(collection)
		return
	}

	// we should never reach here.
	log.Println("no arguments provided")
	os.Exit(1)
}
