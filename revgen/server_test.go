// Copyright 2017 <CompanyName>, Inc. All Rights Reserved.

package revgen

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// Test GET /reverse/response 200
func TestGetSuccess(t *testing.T) {

	// First, need to create a record
	handler := HandleReverseRequest

	form := url.Values{}
	form.Add("phrase", "word1")
	req, err := http.NewRequest("POST", "http://example.com/reverse", strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d\n", w.Code)
	}

	// Poor practice. Use concurrency primitive instead.
	time.Sleep(7 * time.Second)

	uuid := strings.TrimSpace(w.Body.String())
	handler2 := HandleGetResult
	req2, err2 := http.NewRequest("GET", "http://example.com/reverse/"+uuid, nil)
	if err2 != nil {
		log.Fatal(err2)
	}

	w2 := httptest.NewRecorder()
	handler2(w2, req2)
	if w2.Code != 200 {
		t.Errorf("Expected code 200, received code %d\n", w2.Code)
	}
}

// Test GET /reverse/ response 400
func TestGetIllegalId(t *testing.T) {

	handler := HandleGetResult
	req, err := http.NewRequest("GET", "http://example.com/reverse/abcd", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler(w, req)

	fmt.Printf("%d - %s", w.Code, w.Body.String())
	if w.Code != 400 {
		t.Errorf("Expected code 400, received code %d\n", w.Code)
	}
}

// Test GET /reverse/ response 404
func TestGetNotFound(t *testing.T) {

	handler := HandleGetResult
	req, err := http.NewRequest("GET", "http://example.com/reverse/9b53dea0-02a2-4e4f-9337-82651f0bc664", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler(w, req)

	fmt.Printf("%d - %s", w.Code, w.Body.String())
	if w.Code != 404 {
		t.Errorf("Expected 404, got %d\n", w.Code)
	}
}

// Test GET /stats response 200
func TestGetStats(t *testing.T) {

	handler := HandleStatsRequest

	req, err := http.NewRequest("GET", "http://example.com/stats", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler(w, req)

	fmt.Printf("%d - %s", w.Code, w.Body.String())
	if w.Code != 200 {
		t.Errorf("Expected code 200, received code %d\n", w.Code)
	}
}

// Test POST /reverse response 200
func TestPostSuccess(t *testing.T) {

	handler := HandleReverseRequest

	form := url.Values{}
	form.Add("phrase", "word1")
	req, err := http.NewRequest("POST", "http://example.com/reverse", strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler(w, req)

	fmt.Printf("%d - %s", w.Code, w.Body.String())
	if w.Code != 200 {
		t.Errorf("Expected 200, got %d\n", w.Code)
	}
}

// Test POST /reverse response 400
func TestPostNoPhrase(t *testing.T) {

	handler := HandleReverseRequest

	// JSON representation is not supported, so 400 expected
	jsonBytes := []byte(`{"phrase":"word1"}`)
	req, err := http.NewRequest("POST", "http://example.com/reverse", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler(w, req)

	fmt.Printf("%d - %s", w.Code, w.Body.String())
	if w.Code != 400 {
		t.Errorf("Expected 400, got %d\n", w.Code)
	}
}

func TestStartServer(t *testing.T) {

	StartServer("localhost", "8080")
	// serve for less than one second
	time.Sleep(100 * time.Millisecond)
	// send shutdown command
	close(quit)
	time.Sleep(100 * time.Millisecond)
}
