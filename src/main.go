package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
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
func main() {
	http.HandleFunc("/", HomePage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func returnArchiveUrls(domain string) [][]string {
	const base_url string = "https://web.archive.org/cdx/search/cdx?"
	params := url.Values{}
	get_parameters := map[string]string{
		"url":       domain,
		"matchType": "domain",
		"filter":    "statuscode:200",
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

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
		err := tmpl.ExecuteTemplate(w, "base.html", nil)
		checkError(err)
		return
	} else if r.Method == http.MethodPost {
		domain := r.FormValue("domain")
		d := returnArchiveUrls(domain)
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/result.html"))
		err := tmpl.ExecuteTemplate(w, "base.html", d)
		checkError(err)
		return
	}

}
