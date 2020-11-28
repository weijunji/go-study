package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	r := gin.Default()
	g := r.Group("/")
	SetupRouter(g, g)

	testData := []struct {
		body string
		code int
	}{
		{``, http.StatusBadRequest},
		{`{"username": "haruka","password": "123457"}`, http.StatusUnauthorized},
		{`{"username": "haruka","password": "123456"}`, http.StatusOK},
	}

	for _, test := range testData {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(test.body))
		r.ServeHTTP(w, req)

		assert.Equal(t, test.code, w.Code)
	}
}
