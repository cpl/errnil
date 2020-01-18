package badger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cpl.li/go/errnil/pkg/errnil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"cpl.li/go/errnil/pkg/store"
)

func TestBadger(t *testing.T) {
	t.Parallel()

	storage := store.NewPrimitiveStore()
	router := gin.New()
	router.GET("/badge", Badge(storage))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/badge", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, newBadge("missing repo", "critical", "", true).JSON()+"\n", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/badge?repo=cpl.li/go/alpha&style=testing", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, newBadge("nil", "incative", "testing", true).JSON()+"\n", w.Body.String())

	storage.SetEntry("cpl.li/go/alpha", []errnil.Position{
		{
			Filename: "test",
			Line:     1,
			Column:   2,
			Offset:   3,
		},
		{
			Filename: "test",
			Line:     1,
			Column:   2,
			Offset:   3,
		},
	})

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/badge?repo=cpl.li/go/alpha&style=flat-square", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, newBadge("2", badgeColor, "flat-square", false).JSON()+"\n", w.Body.String())
}
