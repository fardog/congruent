// Package example contains an example test case against the public
// mkwords API (https://mkwords.fardog.io/api) and a locally running version.
//
// To test this out, you can pull the docker image locally using:
//
//    docker pull fardog/mkwords
//    docker run -d -p 3000:3000 fardog/mkwords
//
// Then run the tests from this directory using `go test`
package example

import (
	"encoding/json"
	"testing"

	"github.com/fardog/congruent"
)

var servers congruent.Servers

// GenerateResponse represents the structure of a response from the mkwords
// "generate" api endpoint
type GenerateResponse struct {
	Ok             bool     `json:"ok"`
	Result         []string `json:"result"`
	CandidateCount int      `json:"candidate-count"`
}

func init() {
	// define the servers against which we'll be testing; this can be any number
	// of servers, although we only define two here
	servers = []*congruent.Server{
		congruent.NewServer("https://mkwords.fardog.io/api/v1/", nil),
		congruent.NewServer("http://localhost:3000/api/v1/", nil),
	}
}

// TestGenerate tests the simple generate endpoint case, with default values
func TestGenerate(t *testing.T) {
	// create a new request, a GET against the generate endpoint
	request := congruent.NewRequest("GET", "generate", nil, nil)

	// perform the request
	responses, err := servers.Request(request)
	if err != nil {
		t.Error(err)
	}

	// verify the status codes on both responses matched `200`
	if err := responses.StatusEqual(200); err != nil {
		t.Error(err)
	}

	// ensure both responses contain the application/json content header
	if err := responses.HeaderEqual(
		"content-type", "application/json; charset=utf-8"); err != nil {
		t.Error(err)
	}

	// now we'll do some content tests; since mkwords creates a random list of
	// words on each call, we can't just say "are these the same", because the
	// word list will not be. we'll test things individually instead
	for _, resp := range responses {
		// unmarshal the response into our pre-defined struct; any errors here would
		// mean that the response didn't conform to what we expected
		body := GenerateResponse{}
		if err := json.Unmarshal(resp.Body, &body); err != nil {
			t.Error(err)
			t.Fail()
		}

		// verify individual properties in the unmarshal'd response
		if l := len(body.Result); l != 4 {
			t.Errorf("Expected 4 results, got %v", l)
		}

		if !body.Ok {
			t.Errorf("Expected OK to be true, got %v", body.Ok)
		}

		if body.CandidateCount != 70806 {
			t.Errorf(
				"Expected candidate_count to be 70806, got %v",
				body.CandidateCount)
		}
	}
}

// TestGenerateOptions tests the "generate" endpoint with the options it
// supports, asking for parameters that are not the defaults. Otherwise, this
// test is very similar to the TestGenerate above.
func TestGenerateOptions(t *testing.T) {
	request := congruent.NewRequest(
		"GET",
		"generate?min-chars=10&max-chars=20&num-words=10",
		nil,
		nil)

	responses, err := servers.Request(request)
	if err != nil {
		t.Error(err)
	}

	if err := responses.StatusEqual(200); err != nil {
		t.Error(err)
	}

	if err := responses.HeaderEqual(
		"content-type", "application/json; charset=utf-8"); err != nil {
		t.Error(err)
	}

	for _, resp := range responses {
		body := GenerateResponse{}
		if err := json.Unmarshal(resp.Body, &body); err != nil {
			t.Error(err)
			t.Fail()
		}

		if l := len(body.Result); l != 10 {
			t.Errorf("Expected 10 results, got %v", l)
		}

		if !body.Ok {
			t.Errorf("Expected OK to be true, got %v", body.Ok)
		}

		if body.CandidateCount != 21663 {
			t.Errorf(
				"Expected candidate_count to be 21663, got %v",
				body.CandidateCount)
		}
	}
}

// Test404 verifies the 404 response; gets a non-existent endpoint from mkwords
// and verifies it matches expectations.
func Test404(t *testing.T) {
	request := congruent.NewRequest("GET", "bla", nil, nil)

	responses, err := servers.Request(request)
	if err != nil {
		t.Error(err)
	}

	if err := responses.StatusEqual(404); err != nil {
		t.Error(err)
	}
}
