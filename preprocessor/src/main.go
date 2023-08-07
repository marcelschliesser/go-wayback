package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

type Job struct {
	OrderedAt time.Time
	Domain    string
	Email     string
	Year      string
}

var httpClient = &http.Client{}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/user-input.html"))
		err := tmpl.ExecuteTemplate(w, "base.html", nil)
		checkError(err)
	} else if r.Method == http.MethodPost {
		var j Job
		j.Domain = r.FormValue("domain")
		j.Year = r.FormValue("year")
		j.Email = r.FormValue("email")
		j.OrderedAt = time.Now()
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/result.html"))
		err := tmpl.ExecuteTemplate(w, "base.html", j)
		checkError(err)
	}

}

func main() {
	http.HandleFunc("/", indexPage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
