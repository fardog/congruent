package congruent

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestStatus(t *testing.T) {
	gu, err := url.Parse("http://localhost/")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	bu, err := url.Parse("http://bad/")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	mockReq := &http.Request{Method: "GET", URL: gu}
	mockBadReq := &http.Request{Method: "GET", URL: bu}

	responses := Responses{
		&Response{mockReq, nil, nil, 200},
		&Response{mockReq, nil, nil, 200},
		&Response{mockReq, nil, nil, 200}}

	if err := responses.StatusEqual(200); err != nil {
		t.Error(err)
	}
	if err := responses.StatusSame(); err != nil {
		t.Error(err)
	}
	if err := responses.StatusEqual(201); err == nil {
		t.Error("Expected error, but got none!")
	}

	responses = Responses{
		&Response{mockReq, nil, nil, 200},
		&Response{mockBadReq, nil, nil, 201},
		&Response{mockReq, nil, nil, 200}}

	if err := responses.StatusEqual(200); err == nil {
		t.Error("Expected error, but got none!")
	} else {
		if !strings.Contains(err.Error(), "bad") || !strings.Contains(err.Error(), "201") {
			t.Errorf("Did not get expected error string, got: %v", err)
		}
	}
	if err := responses.StatusSame(); err == nil {
		t.Error("Expected error, but got none!")
	} else {
		if !strings.Contains(err.Error(), "bad") || !strings.Contains(err.Error(), "201") {
			t.Errorf("Did not get expected error string, got: %v", err)
		}
	}

}

func TestBodySame(t *testing.T) {
	gu, err := url.Parse("http://localhost/")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	bu, err := url.Parse("http://bad/")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	mockReq := &http.Request{Method: "GET", URL: gu}
	mockBadReq := &http.Request{Method: "GET", URL: bu}

	responses := Responses{
		&Response{mockReq, nil, []byte{'a', 'b', 'c'}, 200},
		&Response{mockReq, nil, []byte{'a', 'b', 'c'}, 200},
		&Response{mockReq, nil, []byte{'a', 'b', 'c'}, 200}}

	if err := responses.BodySame(); err != nil {
		t.Error(err)
	}

	responses = Responses{
		&Response{mockReq, nil, []byte{'a', 'b', 'c'}, 200},
		&Response{mockBadReq, nil, []byte{'a', 'b', 'd'}, 200},
		&Response{mockReq, nil, []byte{'a', 'b', 'c'}, 200}}

	if err := responses.BodySame(); err == nil {
		t.Error("Expected error, but got none!")
	}
}

func TestBodyContentSame(t *testing.T) {
	gu, err := url.Parse("http://localhost/")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	bu, err := url.Parse("http://bad/")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	gs1 := []byte(`{"a":{"b":"c"}}`)
	gs2 := []byte(`
    {
      "a": {
        "b": "c"
      }
    }`)
	bs1 := []byte(`{"a":{"b":"d"}}`)

	mockReq := &http.Request{Method: "GET", URL: gu}
	mockBadReq := &http.Request{Method: "GET", URL: bu}

	responses := Responses{
		&Response{mockReq, nil, gs1, 200},
		&Response{mockReq, nil, gs2, 200},
		&Response{mockReq, nil, gs1, 200}}

	if err := responses.BodyContentSame(); err != nil {
		t.Error(err)
	}

	responses = Responses{
		&Response{mockReq, nil, gs2, 200},
		&Response{mockBadReq, nil, bs1, 200},
		&Response{mockReq, nil, gs2, 200}}

	if err := responses.BodyContentSame(); err == nil {
		t.Error("Expected error, but got none!")
	}
}
