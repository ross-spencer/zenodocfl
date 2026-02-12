/*
lister is responsible for listing urls associated with a collection
in the Mediathek. The output is used as an input for `gather`.

General control flow:

 1. create search url.
 1. request data.
 1. read data and find table.
 1. list urls from table.
 1. given list, provide it as an input to `gather`.
*/
package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ross-spencer/zenodocfl/internal/logformatter"
	"github.com/ross-spencer/zenodocfl/internal/types"
	"golang.org/x/net/html"

	"flag"
)

var (
	search     string
	collection int
	lang       string
	results    int
	checklist  bool
	vers       bool

	// app constants.
	version = "dev-0.0.0"
	commit  = "000000000000000000000000000000000baddeed"
	date    = "1970-01-01T00:00:01Z"
)

var agent string = fmt.Sprintf("INK-lister/%s", version)

const defaultResults int = 10
const defaultLanguage string = "de"

// initFlags initializes the flags we use with this app.
func initFlags() {
	flag.StringVar(&search, "search", "", "string to search for in INK")
	flag.IntVar(&collection, "collection", 0, "collection number to use if known")
	flag.StringVar(&lang, "language", "de", fmt.Sprintf("language, default: '%s'", defaultLanguage))
	flag.IntVar(&results, "results", defaultResults, fmt.Sprintf("number of results to return, default: %d]", defaultResults))
	flag.BoolVar(&checklist, "checklist", false, "output a checklist")
	flag.BoolVar(&vers, "version", false, "return version")
}

// makeResultParams creates a base64 set of result parameters to
// enable download of all results in a single request.
func makeResultParams(number int) string {

	/* Example parameters encoded as JSON:
	   {"from":0,"size":300}
	*/

	log.Printf("maximum number of results requested: %d", number)
	s := fmt.Sprintf("{\"from\": 0, \"size\": %d}", number)
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// extractTable extracts metadata and hyperlinks from Katalog retults.
func extractTable(reader *bufio.Reader, lang string) ([]types.MediathekRecord, error) {
	data, err := html.Parse(reader)
	if err != nil {
		log.Println("err", err)
	}
	m := make([]types.MediathekRecord, 0)
	processTableResults(data, &m, lang)
	if err != nil {
		return []types.MediathekRecord{}, nil
	}
	return m, nil
}

// validateSearch checks a URL and makes sure we are going to be able
// to use it.
func validateSearch(search string) bool {
	return search != ""
}

// makeINKSearchURL creates a URL we can use to return INK results for
// crawling.
func makeINKSearchURL(lang string, search string, results int) string {
	const url = "https://ink.sammlung.cc/table"
	return fmt.Sprintf("%s/%s?search=%s&cursor=%s", url, lang, search, makeResultParams(results))
}

// makeINKSearchURL creates a URL we can use to return INK results for
// crawling.
func makeINKCollectionURL(lang string, collection int, results int) string {
	const url = "https://ink.sammlung.cc/table"
	return fmt.Sprintf("%s/%s?search=&collections=%d&cursor=%s", url, lang, collection, makeResultParams(results))
}

func main() {

	logformatter.Set("lister", true)

	initFlags()
	flag.Parse()

	if vers {
		fmt.Fprintf(os.Stderr, "%s (%s) commit: %s date: %s\n", agent, version, commit, date)
		os.Exit(0)
	} else if flag.NFlag() < 1 || !validateSearch(search) && collection <= 0 {
		fmt.Fprintln(os.Stderr, "Usage:  ")
		fmt.Fprintln(os.Stderr, "        REQUIRED: [-search]  STRING | [-collection]  INT")
		fmt.Fprintln(os.Stderr, "        OPTIONAL: [-lang]    STRING")
		fmt.Fprintln(os.Stderr, "        OPTIONAL: [-results] INTEGER")
		fmt.Fprintln(os.Stderr, "        OPTIONAL: [-checklist] ")
		fmt.Fprintln(os.Stderr, "        OPTIONAL: [-version] ")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Output: [STRING] {result JSON}")
		fmt.Fprintln(os.Stderr, "Output: [STRING] {result checklist}")
		fmt.Fprintf(os.Stderr, "Output: [STRING] {version: '%s'}\n\n", agent)
		flag.Usage()
		os.Exit(0)
	}

	var inkURL string
	if collection == 0 {
		inkURL = makeINKSearchURL(lang, search, results)
	} else {
		inkURL = makeINKCollectionURL(lang, collection, results)
	}

	log.Printf("requesting: %s", inkURL)

	// create a client to set a URL header.
	resp, err := http.Get(inkURL)
	if err != nil {
		log.Println("problem retrieving data from INK:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	content := bufio.NewReader(resp.Body)

	urlList, _ := extractTable(content, lang)

	if checklist {
		fmt.Print("INK Results checklist:\n\n")
		for _, item := range urlList {
			fmt.Println("[ ]", item.Title, ":", item.Url)
		}
		return
	}
	for _, item := range urlList {
		jsonOut, err := json.Marshal(item)
		if err != nil {
			panic("todo")
		}
		fmt.Println(string(jsonOut))
	}

}
