package msclient

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"gopkg.in/resty.v1"
	"net/http"
	"os"
	"testing"
)

type ipifyResp struct{
	IP string `json:"ip"`
}

func TestClientSimple(t *testing.T) {
	if err := os.Setenv("TEST_BASE_URL", "https://api.ipify.org"); err != nil {
		t.Fatal(err)
	}

	c, err := New(MustEnvCfg("TEST_"), zerolog.Nop(), prometheus.NewRegistry(), "test")
	if err != nil {
		t.Fatal(err)
	}

	resp := &ipifyResp{}

	if err = c.Get("/?format=%s", resp, "json"); err != nil {
		t.Fatal(err)
	}

	if resp.IP == "" {
		t.Fatal("empty response")
	}

	resp = &ipifyResp{}

	var result *resty.Response
	result, err = c.NewRequest().SetResult(resp).SetQueryParam("format", "json").Get("/")
	if err != nil {
		t.Fatal(err)
	}
	if s := result.StatusCode(); s != http.StatusOK {
		t.Fatalf("status: %d", s)
	}

	if resp.IP == "" {
		t.Fatal("empty response")
	}
}

