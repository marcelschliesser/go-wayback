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

var userInputTemplate = template.Must(template.ParseFiles("templates/layout.html", "templates/user-input.html"))
var resultTemplate = template.Must(template.ParseFiles("templates/layout.html", "templates/result.html"))

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
		err := userInputTemplate.ExecuteTemplate(w, "user-input.html", nil)
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

		err = resultTemplate.ExecuteTemplate(w, "result.html", j)
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
