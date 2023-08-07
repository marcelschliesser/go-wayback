package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
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
	log.Printf("GET: %v", finalURL)
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

func partitionUrls(records []Record) []Record {
	partitions := make(map[string]Record)

	for _, entry := range records {
		original := entry.Original
		timestamp, err := strconv.Atoi(entry.Timestamp)
		checkError(err)

		if existing, ok := partitions[original]; !ok {
			partitions[original] = entry
		} else {
			existingTimestamp, err := strconv.Atoi(existing.Timestamp)
			checkError(err)

			if timestamp > existingTimestamp {
				partitions[original] = entry
			}
		}
	}

	// Convert the map to a slice
	var result []Record
	for _, v := range partitions {
		result = append(result, v)
	}

	// Sort the slice by Original field
	sort.Slice(result, func(i, j int) bool {
		return result[i].Original < result[j].Original
	})

	return result
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
		err := tmpl.ExecuteTemplate(w, "base.html", nil)
		checkError(err)
		return
	} else if r.Method == http.MethodPost {
		var d Data
		d.Domain = r.FormValue("domain")
		d.Year = r.FormValue("year")
		d.Records = convertResponse(returnArchiveUrls(d.Domain, d.Year))
		d.Records = partitionUrls(d.Records)
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/result.html"))
		err := tmpl.ExecuteTemplate(w, "base.html", d)
		checkError(err)
		return
	}

}

func main() {
	http.HandleFunc("/", indexPage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
