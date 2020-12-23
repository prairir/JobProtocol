package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	globals "github.com/prairir/JobProtocol/Globals"
)

type queueJson struct {
	result string
}

// testing /api/queue route
// make sure the server is running
func TestWebQueue(t *testing.T) {
	port := fmt.Sprintf("%d", globals.ConnPort)
	resp, err := http.Get("http://localhost:" + port)
	if err != nil {
		t.Errorf("/api/queue failed, are you sure the server is running?\nErr: %s", err)
		return
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var data queueJson
	err = decoder.Decode(&data)
	if err != nil {
		t.Errorf("/api/queue failed, are you sure the server is running?\nErr: %s", err)
		return
	}

	if data.result != "" || resp.StatusCode == http.StatusBadRequest {
		t.Logf("/api/queue sucess, got %s", data.result)
		return
	}
	t.Errorf("/api/queue failed, are you sure the server is running?\nErr: no result from server")
	return
}

// testing if the /api/job route works
// make sure the server is running
func TestWebJob(t *testing.T) {
	port := fmt.Sprintf("%d", globals.ConnPort)

	formData := url.Values{
		"job": {"JOB EQN 1+2"},
	}

	resp, err := http.PostForm("http://localhost:"+port, formData)

	if err != nil {
		t.Errorf("/api/job failed, are you sure the server is running?\nErr: %s", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		t.Logf("/api/job sucess, got 200 response")
		return
	}
	t.Errorf("/api/job failed, are you sure the server is running?")
	return
}
