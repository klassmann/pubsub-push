package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"
)

const (
	credentialsVarName string = "GOOGLE_APPLICATION_CREDENTIALS"
	messageMimetype    string = "application/json"
)

type headers []string

func (h *headers) String() string {
	b := strings.Builder{}

	for _, v := range *h {
		b.WriteString(v)
	}

	return b.String()
}

func (h *headers) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func (h *headers) applyHeaders(ht *http.Request) {
	for _, v := range *h {
		parts := strings.Split(v, "=")

		if len(parts) == 2 {
			ht.Header.Set(parts[0], parts[1])
		} else if len(parts) == 1 {
			ht.Header.Set(parts[0], "")
		}
	}
}

type settings struct {
	ProjectID    string
	Subscription string
	Endpoint     string
	Headers      headers
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
	flag.Var(&s.Headers, "header", "A string that represents a Header to be sent. You can use multiple times. Format: key=value")
	flag.Parse()

	if s.ProjectID == "" || s.Endpoint == "" || s.Subscription == "" {
		flag.Usage()
		os.Exit(-1)
	}

	return &s
}

func post(url string, contentType string, body io.Reader, h *headers) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatalf("I was not possible to create a new request: %v\n", err)
		return nil, err
	}
	req.Header.Set("Content-type", contentType)
	h.applyHeaders(req)
	return http.DefaultClient.Do(req)
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
		os.Exit(-1)
	}

	fmt.Printf("Listening subscription %s:\n", settings.Subscription)
	sub := client.Subscription(settings.Subscription)
	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		b, size := encodeMessage(m)
		buff := bytes.NewBuffer(b)

		resp, err := post(settings.Endpoint, messageMimetype, buff, &settings.Headers) //http.Post(settings.Endpoint, messageMimetype, buff)

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
