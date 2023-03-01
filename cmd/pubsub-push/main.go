package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	push "github.com/klassmann/pubsub-push"
)

const (
	credentialsVarName string = "GOOGLE_APPLICATION_CREDENTIALS"
	messageMimetype    string = "application/json"
)

type settings struct {
	ProjectID    string
	Subscription string
	Endpoint     string
	Headers      push.Headers
}

func getArguments() *settings {
	s := settings{}

	flag.StringVar(&s.ProjectID, "project", "", "Google Cloud Project ID")
	flag.StringVar(&s.Subscription, "sub", "", "Subscription name only, without prefix")
	flag.StringVar(&s.Endpoint, "endpoint", "", "Endpoint, format = http://host:port/path")
	flag.Var(&s.Headers, "header", "A string that represents a Header to be sent. You can use multiple times. Format: key=value")
	flag.Parse()

	if s.ProjectID == "" || s.Endpoint == "" || s.Subscription == "" {
		flag.Usage()
		os.Exit(-1)
	}

	return &s
}

func main() {
	settings := getArguments()

	_, b := os.LookupEnv(credentialsVarName)
	if !b {
		fmt.Printf("define %s with the correct credential or the default credential will be used.\n", credentialsVarName)
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, settings.ProjectID)

	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	sub := client.Subscription(settings.Subscription)
	ok, err := sub.Exists(ctx)

	if !ok || err != nil {
		log.Fatalf("failed to connect to subscription: %s", settings.Subscription)
	}

	fmt.Printf("listening subscription %s:\n", sub.String())
	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		b, size := push.EncodeMessage(m)
		buff := bytes.NewBuffer(b)

		resp, err := push.PostMessage(settings.Endpoint, messageMimetype, buff, &settings.Headers)

		if err != nil {
			log.Printf("[ERROR]: %v\n", err)
			m.Nack()
		} else {
			fmt.Printf("[SUCCESS]: message of %d bytes sent to %s: got HTTP %d. Message: ", size, settings.Endpoint, resp.StatusCode)
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
