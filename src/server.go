package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Create a data instance with the message to be displayed in the template.
	data, err := ioutil.ReadFile("data.json")
	checkError(err)
	// Parse the template file(s).
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	// Execute the template with the provided data.
	err2 := tmpl.Execute(w, string(data))
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	log.Println("Start Export")
	export()
	log.Println("Finished Export")
	log.Println("Start Server")
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}
