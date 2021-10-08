//go:build integration
// +build integration

package route

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var previewUrl = flag.String("previewUrl", "", "Workspace preview URL")
var signadotApiKey = flag.String("signadotApiKey", "", "Signadot API key")

func validateInputArgs(t *testing.T) {
	if *previewUrl == "" {
		t.Error("command line arg required for: previewUrl")
	}
	if *signadotApiKey == "" {
		t.Error("command line arg required for: signadotApiKey")
	}
}

func TestRouteAPI(t *testing.T) {
	validateInputArgs(t)
	req, err := http.NewRequest("GET", *previewUrl+"/route?pickup=577,322&dropoff=115,322", nil)
	req.Header.Set("signadot-api-key", *signadotApiKey)

	if err != nil {
		t.Error(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("service call failed. Status code: %s", strconv.Itoa(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("invalid response body. error: %s", err)
	}

	type RouteResponse struct {
		Pickup  string `json:"pickup"`
		Dropoff string `json:"dropoff"`
		ETA     int    `json:"ETA"`
	}

	var p RouteResponse
	err = json.Unmarshal(body, &p)

	if err != nil {
		t.Errorf("unable to unmarshal response body. error: %s", err)
	}

	assert.True(t, p.ETA >= 0, "ETA should be a non-negative value")
}
