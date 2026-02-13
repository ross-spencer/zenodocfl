/*
crater converts a manifest from gather into a RO-CRATE.

 1. read gather manifest.

 2. download all files and resources.

 3. output a ro-crate JSON.

    folder structure:

    ./ro-crate.json
    ./records/
    ...json
    ./media/
    ...bin
    ./poster/
    ...bin
    ./ancillary/   <-- customizable...
    ...bin...
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/ross-spencer/zenodocfl/internal/logformatter"
)

var (
	crate      string
	additional string
	meta       string
	dryrun     bool
	vers       bool

	// app constants.
	version = "dev-0.0.0"
	commit  = "000000000000000000000000000000000baddeed"
	date    = "1970-01-01T00:00:01Z"
)

var agent string = fmt.Sprintf("INK-crater/%s", version)

// initFlags initializes the flags we use with this app.
func initFlags() {
	flag.StringVar(&crate, "crate", "", "collection manifest to convert to RO-CRATE")
	flag.StringVar(&meta, "meta", "", "metadata for the RO-CRATE")
	flag.StringVar(&additional, "additional", "", "change name of ancillary directory")
	flag.BoolVar(&dryrun, "dry-run", false, "peform a dry-run (dont download files)")
	flag.BoolVar(&vers, "version", false, "Return version")
}

// timestamp is a utility function returning a UNIX timestamp for use
// throughout this app.
func timestamp() int64 {
	return time.Now().Unix()
}

// createCrateDir will create a directory within the ro-crate object
// created by this app.
func createCrateDir(dirName string) {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		log.Println("error creating directory:", err)
		os.Exit(1)
	}
}

// getMeta returns a metadata object from a user input.
func getMeta(meta string) userData {
	data, err := os.ReadFile(meta)
	if err != nil {
		log.Println("error reading metadata:", err)
		os.Exit(1)
	}
	var metaData userData
	err = json.Unmarshal(data, &metaData)
	if err != nil {
		log.Println("error reading metadata:", err)
		os.Exit(1)
	}
	return metaData

}

// makeCrrate creates a RO-CRATE from the given collection manifest.
/*
    The collection looks as follows:

	Items      []Item `json:"records"`          <-- redistribute JSON in crate.
	MediaURLs  []string `json:"media_urls"`     <-- download to media.
	PosterURLs []string `json:"poster_urls"`    <-- download to poster.
*/
func makeCrate(manifest string, meta string, dryrun bool) {

	// read the data.
	collection := readManifest(manifest)

	// read metadata.
	userData := getMeta(meta)

	// create global object.
	crateDir := filepath.Join("output", fmt.Sprintf(
		"ro-crate-%s-%d",
		strings.Replace(userData.Name, " ", "-", -1),
		timestamp(),
	),
	)
	recordsDir := filepath.Join(crateDir, "records")
	mediaDir := filepath.Join(crateDir, "media")
	posterDir := filepath.Join(crateDir, "posters")
	anciliaryDir := filepath.Join(crateDir, "anciliary")

	// create directory layout.
	createCrateDir(crateDir)
	createCrateDir(recordsDir)
	createCrateDir(mediaDir)
	createCrateDir(posterDir)
	createCrateDir(anciliaryDir)

	// move records.
	recordParts := moveRecords(collection.Items, recordsDir, "records")
	mediaParts := downloadFile(collection.MediaURLs, mediaDir, "media", dryrun)
	posterParts := downloadFile(collection.PosterURLs, posterDir, "posters", dryrun)

	// get all parts for the manifest.
	allParts := slices.Concat(recordParts, mediaParts, posterParts)

	// summary info.
	log.Println("rocrate parts:", len(allParts))

	userData.parts = allParts

	rocrateData := makeCrateObj(userData)

	data, err := json.MarshalIndent(rocrateData, "", " ")
	if err != nil {
		log.Println("cannot create rocrate JSON:", err)
		os.Exit(1)
	}

	createCrateObj(filepath.Join(crateDir, crateName), string(data))

}

func main() {

	logformatter.Set("crater", true)

	initFlags()
	flag.Parse()

	if vers {
		fmt.Fprintf(os.Stderr, "%s (%s) commit: %s date: %s\n", agent, version, commit, date)
		os.Exit(0)
	} else if flag.NFlag() < 1 || crate == "" && meta == "" {
		fmt.Fprintln(os.Stderr, "Usage:  ")
		fmt.Fprintln(os.Stderr, "        REQUIRED: [-crate]  STRING")
		fmt.Fprintln(os.Stderr, "        REQUIRED: [-meta]  STRING")
		fmt.Fprintln(os.Stderr, "        OPTIONAL: [-additional]  STRING")
		fmt.Fprintln(os.Stderr, "        OPTIONAL: [-dry-run] ")
		fmt.Fprintln(os.Stderr, "        OPTIONAL: [-version] ")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Output: [DIRECTORY] {ro-crate structure")
		fmt.Fprintf(os.Stderr, "Output: [STRING] {version: '%s'}\n\n", agent)
		flag.Usage()
		os.Exit(0)
	}

	if vers {
		fmt.Println(agent)
		return
	}

	if crate != "" {
		makeCrate(crate, meta, dryrun)
		os.Exit(0)
	}

	// we should never reach here.
	log.Println("no arguments provided")
	os.Exit(1)
}
