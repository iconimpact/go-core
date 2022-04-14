package respond

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iconmobile-dev/go-core/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var testdata = map[string]interface{}{"foo": "bar"}

func TestJSON(t *testing.T) {
	l, err := zap.NewDevelopment()
	assert.NoError(t, err)

	// v string
	w := httptest.NewRecorder()

	JSON(w, l, http.StatusOK, "no struct, string")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), `"no struct, string"`)
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")

	// v struct
	w = httptest.NewRecorder()

	JSON(w, l, http.StatusOK, testdata)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), `{"foo":"bar"}`)
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")
}

func TestJSONError(t *testing.T) {
	l, err := zap.NewDevelopment()
	assert.NoError(t, err)

	// non errors
	w := httptest.NewRecorder()

	JSONError(w, l, nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, w.Body.String(), `"Internal Server Error"`)
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")

	// non application error
	w = httptest.NewRecorder()

	err = fmt.Errorf("some basic error")

	JSONError(w, l, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, w.Body.String(), `"Internal Server Error"`)
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")

	// application error
	w = httptest.NewRecorder()

	JSONError(w, l, errors.E(err, errors.NotFound, "Data not found"))
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, w.Body.String(), `"Data not found"`)
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")
}
