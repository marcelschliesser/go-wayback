package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"

	"cloud.google.com/go/pubsub"
)

type Job struct {
	OrderedAt time.Time
	Domain    string
	Email     string
	Year      string
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	parsedTemplate, _ := template.ParseFiles(path.Join("./templates", tmpl), "templates/layout.html")
	err := parsedTemplate.Execute(w, data)
	checkError(err)
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		renderTemplate(w, "user-input.html", nil)

	} else if r.Method == http.MethodPost {

		var j Job
		j.Domain = r.FormValue("domain")
		j.Year = r.FormValue("year")
		j.Email = r.FormValue("email")
		j.OrderedAt = time.Now()

		PublishMessage("test-topic", j)

		renderTemplate(w, "result.html", j)
	}

}

func PublishMessage(topicName string, data Job) {
	ctx := context.Background()

	jsonData, err := json.Marshal(data)
	checkError(err)

	client, err := pubsub.NewClient(ctx, "go-wayback")
	checkError(err)
	defer client.Close()

	topic := client.Topic(topicName)

	msg := &pubsub.Message{
		Data: jsonData,
	}

	msgID, err := topic.Publish(ctx, msg).Get(ctx)
	checkError(err)

	log.Printf("Message published with ID: %s\n", msgID)
}

func main() {
	http.HandleFunc("/", indexPage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
