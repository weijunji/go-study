package lottery

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetLotteryList(t *testing.T) {
	r := gin.Default()
	g := r.Group("/")
	SetupRouter(g)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/lottery/list", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "5fcf2625436da21d6dca2981")
}

func TestGetLotteryInfo(t *testing.T) {
	r := gin.Default()
	g := r.Group("/")
	SetupRouter(g)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/lottery/info/5fcf2625436da21d6dca2981", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "award1")

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/lottery/info/5fcf2625436da21d6dca2000", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/lottery/info/5fcf2625436da21d6dca2", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
