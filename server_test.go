package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestPutBPTree(t *testing.T) {
	testData := []struct {
		body string
		code int
	}{
		{`{"key": 123, "val": "test0"}`, http.StatusOK},
		{`{"key": 166, "val": "test0"}`, http.StatusOK},
		{`{"key": 123, "val": "test1"}`, http.StatusOK},
		{`{"key": 161, "val": "test0"}`, http.StatusOK},
		{`{"key": 161, "val": "test1"}`, http.StatusOK},
		{`{"key": "ada", "val": "test"}`, http.StatusBadRequest},
		{``, http.StatusBadRequest},
		{`{"val": "test"}`, http.StatusBadRequest},
		{`{"key": "ada"}`, http.StatusBadRequest},
	}
	router := setupRouter()

	for _, test := range testData {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/bptree/tree", strings.NewReader(test.body))
		router.ServeHTTP(w, req)
		assert.Equal(t, test.code, w.Code)
	}
	testData = []struct {
		body string
		code int
	}{
		{`{"key": 123, "val": " "}`, http.StatusOK},
		{`{"key": 166, "val": " "}`, http.StatusOK},
	}
	for _, test := range testData {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("DELETE", "/bptree/tree", strings.NewReader(test.body))
		router.ServeHTTP(w, req)
		assert.Equal(t, test.code, w.Code)
	}
}
func TestDeleteBPTree(t *testing.T) {
	testData := []struct {
		body string
		code int
	}{
		{`{"key": 123, "val": " "}`, http.StatusBadRequest},
		{`{"key": 166, "val": " "}`, http.StatusBadRequest},
	}
	router := setupRouter()

	for _, test := range testData {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/bptree/tree", strings.NewReader(test.body))
		router.ServeHTTP(w, req)
		assert.Equal(t, test.code, w.Code)
	}
}
