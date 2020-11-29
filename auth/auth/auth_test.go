package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestEncrypter(t *testing.T) {
	assert.Equal(t, "0d3022fabc99ac91919a06609b205ebe", encryptPassword("123456"))
}

func TestLogin(t *testing.T) {
	r := gin.Default()
	g := r.Group("/")
	SetupRouter(g, g)

	testData := []struct {
		body string
		code int
	}{
		{``, http.StatusBadRequest},
		{`{"username": "admin","password": "123456"}`, http.StatusUnauthorized},
		{`{"username": "haruka","password": "123457"}`, http.StatusUnauthorized},
		{`{"username": "kanbusi","password": "123457"}`, http.StatusOK},
	}

	for _, test := range testData {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(test.body))
		r.ServeHTTP(w, req)

		assert.Equal(t, test.code, w.Code)
	}
	//w := httptest.NewRecorder()
	//req, _ := http.NewRequest("POST", "/auth/register", strings.NewReader(`{"username": "kanbusi","password1": "123457","password2": "123457"}`))
	//r.ServeHTTP(w, req)
	//
	//assert.Equal(t, http.StatusOK, w.Code)
}
