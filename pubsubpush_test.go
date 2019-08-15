package pubsubpush

import "testing"

func TestHeadersSet(t *testing.T) {
	h := Headers{}
	h.Set("Content-type=application/json")
	h.Set("Auth=api-key")

	if len(h) != 2 {
		t.Errorf("Len expected 2 and got %d", len(h))
	}

	if h[0] != "Content-type=application/json" {
		t.Errorf("Expected value %s and got %s", "Content-type=application/json", h[0])
	}

	if h[1] != "Auth=api-key" {
		t.Errorf("Expected value %s and got %s", "Auth=api-key", h[1])
	}
}

func TestHeadersString(t *testing.T) {
	h := Headers{}
	h.Set("Content-type=application/json")

	if len(h) != 1 {
		t.Errorf("Len expected 2 and got %d", len(h))
	}

	if h.String() != "Content-type=application/json" {
		t.Errorf("Expected value %s and got %s", "Content-type=application/json", h.String())
	}

}
