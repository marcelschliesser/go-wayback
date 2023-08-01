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

type Record struct {
	URLKey     string
	Timestamp  string
	Original   string
	MimeType   string
	StatusCode string
	Digest     string
	Length     string
}

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
func returnArchiveUrls(domain, year string) [][]string {
	const base_url string = "https://web.archive.org/cdx/search/cdx?"
	params := url.Values{}
	get_parameters := map[string]string{
		"url":       domain,
		"matchType": "domain",
		"filter":    "statuscode:200",
		"output":    "json",
		"from":      year,
		"to":        year}

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

func convertResponse(input [][]string) []Record {
	var records []Record
	for _, values := range input[1:] {
		records = append(records, Record{
			URLKey:     values[0],
			Timestamp:  values[1],
			Original:   values[2],
			MimeType:   values[3],
			StatusCode: values[4],
			Digest:     values[5],
			Length:     values[6],
		})
	}

	return records
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
		err := tmpl.ExecuteTemplate(w, "base.html", nil)
		checkError(err)
		return
	} else if r.Method == http.MethodPost {
		domain := r.FormValue("domain")
		year := r.FormValue("year")
		d := returnArchiveUrls(domain, year)
		converted := convertResponse(d)
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/result.html"))
		err := tmpl.ExecuteTemplate(w, "base.html", converted)
		checkError(err)
		return
	}

}
