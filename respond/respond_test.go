package respond

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iconmobile-dev/go-core/errors"
	"github.com/stretchr/testify/assert"
)

var testdata = map[string]interface{}{"foo": "bar"}

func TestJSON(t *testing.T) {
	// v string
	r := httptest.NewRequest("GET", "http://example.com/v1", nil)
	w := httptest.NewRecorder()

	JSON(w, r, http.StatusOK, "no struct, string")

	assert.Equal(t, http.StatusOK, w.Code)
	var data JSONMsg
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &data))
	assert.Equal(t, data.Data, nil)
	assert.Equal(t, data.Msg, "no struct, string")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")

	// v struct
	r = httptest.NewRequest("GET", "http://example.com/v1", nil)
	w = httptest.NewRecorder()

	JSON(w, r, http.StatusOK, &testdata)

	assert.Equal(t, http.StatusOK, w.Code)
	data = JSONMsg{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &data))
	assert.Equal(t, data.Data, testdata)
	assert.Equal(t, data.Msg, "")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")
}

func TestJSONError(t *testing.T) {
	// non application error
	r := httptest.NewRequest("GET", "http://example.com/v1", nil)
	w := httptest.NewRecorder()

	err := fmt.Errorf("some basic error")

	JSONError(w, r, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var data JSONMsg
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &data))
	assert.Equal(t, data.Data, nil)
	assert.Equal(t, data.Msg, "Internal Server Error")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")

	// application error
	r = httptest.NewRequest("GET", "http://example.com/v1", nil)
	w = httptest.NewRecorder()

	JSONError(w, r, errors.E(err, errors.NotFound, "Data not found"))
	data = JSONMsg{}
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &data))
	assert.Equal(t, data.Data, nil)
	assert.Equal(t, data.Msg, "Data not found")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")
}
