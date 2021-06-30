package msclient

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"os"
	"testing"
)

func TestClientSimple(t *testing.T) {
	if err := os.Setenv("TEST_BASE_URL", "https://api.ipify.org"); err != nil {
		t.Fatal(err)
	}

	c, err := New(MustEnvCfg("TEST_"), zerolog.Nop(), prometheus.NewRegistry(), "test")
	if err != nil {
		t.Fatal(err)
	}

	resp := &struct{
		IP string `json:"ip"`
	}{}

	if err = c.get("/?format=%s", resp, "json"); err != nil {
		t.Fatal(err)
	}

	if resp.IP == "" {
		t.Fatal("empty response")
	}
}
