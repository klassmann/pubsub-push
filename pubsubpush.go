package pubsubpush

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/pubsub"
)

type message struct {
	MessageID  string            `json:"messageId"`
	Data       string            `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

type request struct {
	Message message `json:"message"`
}

// Headers is the list of HTTP Header to be applied on Request
type Headers []string

func (h *Headers) String() string {
	b := strings.Builder{}

	for _, v := range *h {
		b.WriteString(v)
	}

	return b.String()
}

// Set appends a new Header
func (h *Headers) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func (h *Headers) applyHeaders(ht *http.Request) {
	for _, v := range *h {
		parts := strings.Split(v, "=")

		if len(parts) == 2 {
			ht.Header.Set(parts[0], parts[1])
		} else if len(parts) == 1 {
			ht.Header.Set(parts[0], "")
		}
	}
}

// EncodeMessage prepares the message to be like the HTTP Push from PubSub
// It is a JSON with a data field containing a base64 value
func EncodeMessage(m *pubsub.Message) ([]byte, int) {
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

// PostMessage sends the message to endpoint
func PostMessage(url string, contentType string, body io.Reader, h *Headers) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatalf("I was not possible to create a new request: %v\n", err)
		return nil, err
	}
	req.Header.Set("Content-type", contentType)
	h.applyHeaders(req)
	return http.DefaultClient.Do(req)
}
