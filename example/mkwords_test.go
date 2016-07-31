package example

import (
	"encoding/json"
	"testing"

	"github.com/fardog/congruent"
)

var servers congruent.Servers

type GenerateResponse struct {
	Ok             bool     `json:"ok"`
	Result         []string `json:"result"`
	CandidateCount int      `json:"candidate-count"`
}

func init() {
	servers = []*congruent.Server{
		congruent.NewServer("https://mkwords.fardog.io/api/v1/", nil),
		congruent.NewServer("http://localhost:3000/api/v1/", nil),
	}
}

func TestGenerate(t *testing.T) {
	request := congruent.NewRequest("GET", "generate", nil, nil)

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
