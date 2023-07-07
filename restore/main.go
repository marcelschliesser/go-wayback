package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

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
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // return []byte

	var result [][]string

	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		log.Fatal(err2)
	}
	log.Printf("Wayback URL Count: %d", len(result))
	return result
}

func main() {
	log.SetFlags(log.LstdFlags) // Include timestamp in log messages

	urls := returnArchiveUrls("xxl-angeln.de")

	for _, v := range urls {
		fmt.Printf("http://web.archive.org/web/%vif_/%v", v[1], v[2])
		fmt.Println("")
	}

}
