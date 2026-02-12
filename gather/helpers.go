package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ross-spencer/zenodocfl/internal/types"
)

// printManifest is a helpter to output JSON to stdout.
//
//lint:ignore U1000 this is a convenience function.
func printManifest(manifest []inkRecord) {
	jsonOut, err := json.Marshal(manifest)
	if err != nil {
		log.Println("problem outputting JSON record:", err)
	}
	fmt.Println(string(jsonOut))
}

// printCollection is a helpter to output JSON to stdout. If an
// output string is supplied it will write it to a collections file.
func printCollection(collection types.Collection, output string) {
	jsonOut, err := json.MarshalIndent(collection, "", " ")
	if err != nil {
		log.Println("problem outputting JSON record:", err)
	}
	if output != "" {
		collectionFile, err := os.Create(fmt.Sprintf("%s.collection", output))
		if err != nil {
			log.Println("problem creating collection file:", err)
		}
		defer collectionFile.Close()
		collectionFile.WriteString(string(jsonOut))
		collectionFile.WriteString("\n")
		return
	}
	fmt.Println(string(jsonOut))
}
