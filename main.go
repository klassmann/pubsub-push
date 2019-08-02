package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
)

const (
	credentialsVarName string = "GOOGLE_APPLICATION_CREDENTIALS"
	messageMimetype    string = "application/json"
)

type settings struct {
	ProjectID    string
	Subscription string
	Endpoint     string
}

type message struct {
	MessageID  string            `json:"messageId"`
	Data       string            `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

type request struct {
	Message message `json:"message"`
}

func encodeMessage(m *pubsub.Message) ([]byte, int) {
	data := m.Data
	req := request{}
	req.Message.Data = base64.StdEncoding.EncodeToString(data)
	req.Message.Attributes = m.Attributes
	req.Message.MessageID = m.ID
	b, err := json.Marshal(req)

	if err != nil {
		log.Fatal(err)
		return []byte{}, 0
	}

	return b, len(b)
}

func getArguments() *settings {
	s := settings{}

	flag.StringVar(&s.ProjectID, "project", "", "Google Cloud Project ID")
	flag.StringVar(&s.Subscription, "sub", "", "Subscription name only, without prefix")
	flag.StringVar(&s.Endpoint, "endpoint", "", "Endpoint, format = http://host:port/path")
	flag.Parse()

	if s.ProjectID == "" || s.Endpoint == "" || s.Subscription == "" {
		flag.Usage()
	}

	return &s
}

func main() {
	settings := getArguments()

	_, b := os.LookupEnv(credentialsVarName)
	if !b {
		fmt.Printf("You need to define %s variable with the correct credentials.\n", credentialsVarName)
		os.Exit(-1)
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, settings.ProjectID)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Printf("Listening subscription %s:\n", settings.Subscription)
	sub := client.Subscription(settings.Subscription)
	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		b, size := encodeMessage(m)
		buff := bytes.NewBuffer(b)

		resp, err := http.Post(settings.Endpoint, messageMimetype, buff)

		if err != nil {
			log.Fatalf("Error on send message to endpoint: %v\n", err)
			m.Nack()
		} else {
			fmt.Printf("Message with %d bytes was sent to %s, got HTTP %d. Message: ", size, settings.Endpoint, resp.StatusCode)
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				fmt.Printf("ACK\n")
				m.Ack()
			} else {
				fmt.Printf("NACK\n")
				m.Nack()
			}
		}
	})
}
