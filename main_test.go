package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func handlerFunc(f func(w http.ResponseWriter, r *http.Request)) (string, error) {
	ts := httptest.NewServer(http.HandlerFunc(f))
	defer ts.Close()
	res, err := http.Get(ts.URL)

	if err != nil {
		return "", fmt.Errorf(("err %s"), err)
	}

	rbody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return "", fmt.Errorf(("err %s"), err)
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf(("Status code error: %d %s"), res.StatusCode, res.Status)
	}

	return string(rbody), nil
}

func TestRootHandler(t *testing.T) {
	body, err := handlerFunc(RootHandler)
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if body != "Hello, World!" {
		t.Errorf("Expected %s, got %s", "Hello, World!", body)
	}
}

func TestHealthCheckHandler(t *testing.T) {
	body, err := handlerFunc(HealthCheckHandler)
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if body != "ok" {
		t.Errorf("Expected %s, got %s", "ok", body)
	}
}
