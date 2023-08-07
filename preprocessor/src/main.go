package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %v", err)
	}

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return "", fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	topic := client.Topic(topicName)

	msg := &pubsub.Message{
		Data: jsonData,
	}

	msgID, err := topic.Publish(ctx, msg).Get(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to publish message: %v", err)
	}

	return msgID, nil
}

func main() {
	http.HandleFunc("/", indexPage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
