package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

type queueJSON struct {
	queue []interface{}
}

// testing /api/queue route
// make sure the server is running
func TestWebQueue(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/api/queue")
	if err != nil {
		t.Errorf("/api/queue failed, are you sure the server is running?\nErr: %s", err)
		return
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var data queueJSON
	err = decoder.Decode(&data)
	if err != nil {
		t.Errorf("/api/queue failed, are you sure the server is running?\nErr: %s", err)
		return
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Logf("/api/queue sucess, got %s", data.queue)
		return
	}
	t.Errorf("/api/queue failed, are you sure the server is running?\nErr: no result from server")
	return
}

// testing if the /api/job route works
// make sure the server is running
func TestWebJob(t *testing.T) {
	formData := url.Values{
		"job": {"JOB EQN 1+2"},
	}

	resp, err := http.PostForm("http://localhost:8080/api/job", formData)

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
