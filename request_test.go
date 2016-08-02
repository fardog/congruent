package congruent

import (
	"testing"
)

func TestPrepareBodyString(t *testing.T) {
	expected := "heyooo"

	req := &Request{Body: expected}
	buf, err := req.PrepareBody()
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if string(buf.Bytes()) != expected {
		t.Errorf("expected %s, got %s", expected, buf.Bytes())
	}
}

func TestPrepareBodyJSON(t *testing.T) {
	expected := `{"ok":true,"items":["a","b","c"]}`

	type testData struct {
		Ok    bool     `json:"ok"`
		Items []string `json:"items"`
	}

	req := &Request{Body: &testData{true, []string{"a", "b", "c"}}}
	buf, err := req.PrepareBody()
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if string(buf.Bytes()) != expected {
		t.Errorf("expected %s, got %s", expected, buf.Bytes())
	}
}
