package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// HTTP client that can be reused for multiple requests.
var httpClient = &http.Client{}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func returnArchiveUrls(domain string) [][]string {
	const base_url string = "https://web.archive.org/cdx/search/cdx?"
	params := url.Values{}
	get_parameters := map[string]string{
		"url":       domain,
		"matchType": "domain",
		"filtes":    "statuscode:200",
		"output":    "json",
		"from":      "2011",
		"to":        "2013"}

	for k, v := range get_parameters {
		params.Add(k, v)
	}
	url := fmt.Sprintf(base_url + params.Encode())
	log.Printf("Base Url: %v", url)
	resp, err := httpClient.Get(url)
	checkError(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // return []byte
	checkError(err)

	var result [][]string

	err = json.Unmarshal(body, &result)
	checkError(err)

	return result
}

func export() {
	// Include timestamp in log messages
	log.SetFlags(log.LstdFlags)

	urls := returnArchiveUrls("xxl-angeln.de")

	fieldNames := urls[0]
	records := make([]map[string]string, len(urls)-1)
	for i, record := range urls[1:] {
		recordMap := make(map[string]string, len(fieldNames))

		for j, field := range record {
			recordMap[fieldNames[j]] = field
		}

		records[i] = recordMap
	}
	data, err := json.Marshal(records)
	checkError(err)
	ioutil.WriteFile("data.json", data, 0644)
	log.Printf("Wayback URL Count: %d", len(records))
	log.Printf("Wayback URL Count: %d", len(urls))
}
