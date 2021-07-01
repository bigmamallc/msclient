package mocksvc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestMockServiceSimple(t *testing.T) {
	type Json struct {
		Foo string `json:"foo"`
	}

	s := Start(func(req *http.Request) (int, interface{}) {
		return 200, &Json{
			Foo: "bar",
		}
	})

	resp, err := http.Get(s.BaseURL())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("status code: %d", resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	j := &Json{}
	if err := json.Unmarshal(data, j); err != nil {
		t.Fatal(err)
	}
	if j.Foo != "bar" {
		t.Fatalf("unexpected: %s", j.Foo)
	}

	s.Close()
}
