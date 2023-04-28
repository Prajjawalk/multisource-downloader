package utils

import (
	"errors"
	"log"
	"net/http"
	"strconv"
)

// confirmSupportAndFileChunkSize tests to see if "Accept-Ranges" is part of the HTTP Response header
// If HTTP Range requests are not supported, return server not supported error
// If supported, return the filesize and anticipated chunkSize
func ConfirmSupportAndFileSize(dwLink string) (int64, error) {
	// Set DisableCompression to true (default is false)
	// This ensures Go's internal transport behavior does not mess with our logic
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	response, err := client.Get(dwLink)
	if err != nil {
		log.Fatalln(err)
		return 0, errors.New("HTTP error: GET request failed")
	}
	acceptRanges := response.Header["Accept-Ranges"]
	if len(acceptRanges) == 0 || acceptRanges[0] == "none" {
		return 0, errors.New("server error: Accept-Ranges Header does not exist in HTTP Response")
	}
	filesize, err := strconv.ParseInt(response.Header["Content-Length"][0], 10, 64)
	return filesize, err
}
