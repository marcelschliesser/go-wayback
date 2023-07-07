package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	log.SetFlags(log.LstdFlags) // Include timestamp in log messages

	var base_url string = "https://web.archive.org/cdx/search/cdx?"

	params := url.Values{}

	get_parameters := map[string]string{
		"url":       "xxl-angeln.de",
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

	log.Printf("Body Type is: %T", body)

	var result [][]string

	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		log.Fatal(err2)
	}
	log.Printf("Wayback URL Count: %d", len(result))
	for _, v := range result {
		fmt.Printf("http://web.archive.org/web/%vif_/%v", v[1], v[2])
		fmt.Println("")
	}

}
