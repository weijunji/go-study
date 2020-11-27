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

func TestBPTree(t *testing.T) {
	router := setupRouter()
	// Test put
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
	for _, test := range testData {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/bptree/tree", strings.NewReader(test.body))
		router.ServeHTTP(w, req)
		assert.Equal(t, test.code, w.Code)
	}

	// test delete
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

	// test get
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bptree/tree", strings.NewReader(`{"key": 161}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"value":"test1"}`, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/bptree/tree", strings.NewReader(`{"key": 123313}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteBPTreeError(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/bptree/dtree", strings.NewReader(`{"key": 250}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/bptree/niltree", strings.NewReader(`{"key": 250}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/bptree/dtree", strings.NewReader(`{"key": 250, "val": "test1"}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/bptree/dtree", strings.NewReader(`{"key": 251}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
