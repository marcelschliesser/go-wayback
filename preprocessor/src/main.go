package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
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

		msgID, err := PublishMessage("go-wayback", "test-topic", j)
		checkError(err)
		log.Printf("Message published with ID: %s\n", msgID)

		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/result.html"))
		err = tmpl.ExecuteTemplate(w, "base.html", j)
		checkError(err)
	}

}

// PublishMessage sends a message to a specified Google Cloud Pub/Sub topic.
func PublishMessage(projectID, topicName string, data Job) (string, error) {
	ctx := context.Background()

	jsonData, err := json.Marshal(data)
	checkError(err)

	client, err := pubsub.NewClient(ctx, projectID)
	checkError(err)
	defer client.Close()

	topic := client.Topic(topicName)

	msg := &pubsub.Message{
		Data: jsonData,
	}

	msgID, err := topic.Publish(ctx, msg).Get(ctx)
	checkError(err)

	return msgID, nil
}

func main() {
	http.HandleFunc("/", indexPage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
