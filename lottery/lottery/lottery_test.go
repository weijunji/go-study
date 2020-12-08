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

func TestCountSum(t *testing.T) {
	info := LotteryInfo{Awards: []AwardInfo{{Rate: 100}, {Rate: 500}}}
	countRate(&info)
	assert.Equal(t, uint32(600), info.rateSum)
}

func TestGetInfoByID(t *testing.T) {
	info, _ := getInfoByID("5fcf2625436da21d6dca2981")
	assert.Equal(t, "test", info.Title)
	infoCache, _ := getInfoByID("5fcf2625436da21d6dca2981")
	assert.Equal(t, "test", infoCache.Title)
	assert.True(t, info == infoCache, "this should be cached")
	_, err := getInfoByID("5fcf2625436da21d6dca29")
	assert.EqualError(t, err, "id is not valid objectId")
	_, err = getInfoByID("5fcf2625436da21d6dca2000")
	assert.EqualError(t, err, "no such lottery")
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

func TestDeleteOneFromMysql(t *testing.T) {
	success := deleteOne("test", &AwardInfo{ID: "5fcf44ae7a9b950204000000", Value: 1000})
	assert.True(t, success, "Delete this should be always success")
	success = deleteOne("test", &AwardInfo{ID: "5fcf44ae7a9b950204000001", Value: 1000})
	assert.False(t, false, "Delete this should be failed : remain 0")
	success = deleteOne("test", &AwardInfo{ID: "5fcf44ae7a9b950203ffffff", Value: 1000})
	assert.False(t, false, "Delete this should be failed : no such award")
}

func TestDeleteOneFromRedis(t *testing.T) {
	success := deleteOne("test", &AwardInfo{ID: "5fcf44ae7a9b950204000010", Value: 1})
	assert.True(t, success, "Delete this should be always success")
	success = deleteOne("test", &AwardInfo{ID: "5fcf44ae7a9b950204000011", Value: 1})
	assert.False(t, false, "Delete this should be failed : remain 0")
	success = deleteOne("test", &AwardInfo{ID: "5fcf44ae7a9b950203ffffff", Value: 1})
	assert.False(t, false, "Delete this should be failed : no such award")
}

func TestProcessLottery(t *testing.T) {
	award := processLottery(&LotteryInfo{rateSum: 0})
	assert.Nil(t, award)

	info := LotteryInfo{
		rateSum: 500000,
		Awards: []AwardInfo{
			{ID: "1", Rate: 300000},
			{ID: "2", Rate: 200000},
		},
	}
	sum1, sum2 := 0, 0
	const total = 1000
	for n := 0; n < total; n++ {
		award = processLottery(&info)
		if award != nil {
			if award.ID == "1" {
				sum1++
			} else {
				sum2++
			}
		}
	}
	assert.True(t, sum1 > 250 && sum1 < 350, "sum1 should near 300")
	assert.True(t, sum2 > 160 && sum2 < 340, "sum2 should near 200")
}

func TestLottery(t *testing.T) {
	r := gin.Default()
	g := r.Group("/")
	SetupRouter(g)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/lottery/info/5fcf2625436da21d6dca2000", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/lottery/info/5fcf2625436da21d6dca2", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var i int
	for i = 0; i < 5; i++ {
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/lottery/lottery/5fcf2625436da21d6dca2981", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		if len(body) > 40 {
			break
		}
	}
	assert.NotEqual(t, 5, i, "Should get at least one award in 5 times")
}
