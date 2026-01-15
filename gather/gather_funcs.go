package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

// exists provides a test to check whether a given directory exists.
func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
		log.Println("another error has occurred:", err)
		return false
	}
	return true
}

// readJSON from a file and return a data structure that can be used
// throughout the gather process.
func readJSON(filename string) (inkRecord, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return inkRecord{}, err
	}
	var record inkRecord
	json.Unmarshal(data, &record)
	return record, nil
}

// convertMediaServerURI converts a media server URI to a functioning
// url.
func convertMediaServerURI(url string) string {
	if url == "" {
		return ""
	}
	const oldPrefix string = "mediaserver:hsm/"
	const newPrefix string = "https://ba14ns21403-sec1.fhnw.ch/mediasrv/hsm/"
	return fmt.Sprintf("%s/master", strings.Replace(url, oldPrefix, newPrefix, 1))
}

// convertIdentifier creates a URL from a given identivier.
func convertIdentifier(key string, value string) string {
	/*
	   ark: 	ark:/15737/p658-sjm6-66z4 == https://n2t.net/ark:/15737/p658-sjm6-66z4
	   handle: 	20.500.11806/med/3jzx-tf3s-g1 == https://hdl.handle.net/20.500.11806/med/3jzx-tf3s-g1
	*/

	switch key {
	case "ark":
		return fmt.Sprintf("https://n2t.net/%s", value)
	case "handle":
		return fmt.Sprintf("https://hdl.handle.net/%s", value)
	default:
		log.Println("unknown identifier type:", key)
	}
	return value
}
