package main

import (
	"encoding/json"
	"fmt"

	"github.com/ross-spencer/zenodocfl/internal/types"
)

// printManifest is a helpter to output JSON to stdout.
func printManifest(manifest []inkRecord) {
	jsonOut, err := json.Marshal(manifest)
	if err != nil {
		panic("todo")
	}
	fmt.Println(string(jsonOut))
}

// printCollection is a helpter to output JSON to stdout.
func printCollection(collection types.Collection) {
	jsonOut, err := json.Marshal(collection)
	if err != nil {
		panic("todo")
	}
	fmt.Println(string(jsonOut))
}
