package example

import (
	"testing"

	"github.com/fardog/congruent"
)

var servers congruent.Servers

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
