package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Data struct {
	Domain  string
	Year    string
	Records []Record
}

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
	base_url, _ := url.Parse("https://web.archive.org/cdx/search/cdx?")
	query := url.Values{}
	get_parameters := map[string]string{
		"url":       domain,
		"matchType": "domain",
		"filter":    "statuscode:200",
		"output":    "json",
		"from":      year,
		"to":        year}

	for key, value := range get_parameters {
		query.Set(key, value)
	}
	base_url.RawQuery = query.Encode()
	finalURL := base_url.String()
	log.Printf("Base Url: %v", finalURL)
	resp, err := httpClient.Get(finalURL)
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
	if len(input) > 1 {
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
	} else {
		return nil
	}
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
		var d Data
		d.Domain = domain
		d.Year = year
		d.Records = convertResponse(returnArchiveUrls(domain, year))
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/result.html"))
		err := tmpl.ExecuteTemplate(w, "base.html", d)
		checkError(err)
		return
	}

}
