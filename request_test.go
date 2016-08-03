package congruent

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
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

func TestRequest(t *testing.T) {
	// channel for storing "has been called" ints from test servers
	c := make(chan int, 2)

	ts0 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "a")
		// check header
		if h := r.Header.Get("X-Test"); h != "zero" {
			t.Errorf(`Expected header value "zero", got %s`, h)
		}
		// write to channel to track request
		c <- 0
	}))
	defer ts0.Close()
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "b")
		// check header
		if h := r.Header.Get("X-Test"); h != "one" {
			t.Errorf(`Expected header value "one", got %s`, h)
		}
		// write to channel to track request
		c <- 1
	}))
	defer ts1.Close()

	servers := Servers([]*Server{
		NewServer(ts0.URL, &http.Header{"X-Test": []string{"zero"}}),
		NewServer(ts1.URL, &http.Header{"X-Test": []string{"one"}})})

	request := NewRequest("GET", "/", nil, nil)

	responses, err := servers.Request(request)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	// block until both requests finish, and save results
	var seen []int
	for i := 0; i < 2; i++ {
		seen = append(seen, <-c)
	}

	// verify that we saw responses from both servers
	sort.Sort(sort.IntSlice(seen))
	for i, n := range seen {
		if i != n {
			t.Errorf("Expected server %d to be called, saw %d", i, n)
		}
	}

	// verify we saw responses from each test server
	var resp []string
	for _, response := range responses {
		resp = append(resp, string(response.Body))
	}

	// verify that we saw responses from both servers
	sort.Sort(sort.StringSlice(resp))
	if resp[0] != "a" {
		t.Errorf(`Expected "a", got %s`, resp[0])
	}
	if resp[1] != "b" {
		t.Errorf(`Expected "b", got %s`, resp[1])
	}
}
