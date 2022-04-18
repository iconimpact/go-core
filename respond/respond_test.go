package respond

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iconimpact/go-core/errors"
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
	assert.Equal(t, `"no struct, string"`, w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.HeaderMap.Get("Content-Type"))

	// v struct
	w = httptest.NewRecorder()

	JSON(w, l, http.StatusOK, testdata)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"foo":"bar"}`, w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.HeaderMap.Get("Content-Type"))
}

func TestJSONError(t *testing.T) {
	l, err := zap.NewDevelopment()
	assert.NoError(t, err)

	// non errors
	w := httptest.NewRecorder()

	JSONError(w, l, nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"msg":"Internal Server Error"}`, w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.HeaderMap.Get("Content-Type"))

	// non application error
	w = httptest.NewRecorder()

	err = fmt.Errorf("some basic error")

	JSONError(w, l, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"msg":"Internal Server Error"}`, w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.HeaderMap.Get("Content-Type"))

	// application error
	w = httptest.NewRecorder()

	JSONError(w, l, errors.E(err, errors.NotFound, "Data not found"))

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, `{"msg":"Data not found"}`, w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.HeaderMap.Get("Content-Type"))

	// custom error response
	errorRsp := func(err error) interface{} {
		var status int
		var errMsg string

		// set custom app err Message
		appErr, ok := err.(*errors.Error)
		if !ok {
			status = http.StatusInternalServerError
			errMsg = "Internal Server Error"
		} else {
			status = errors.ToHTTPStatus(appErr)
			errMsg = errors.ToHTTPResponse(appErr)
		}

		rsp := struct {
			Msg    string `json:"msg"`
			Status int    `json:"status"`
		}{
			Msg:    errMsg,
			Status: status,
		}

		return rsp
	}
	SetJSONErrorResponse(errorRsp)

	w = httptest.NewRecorder()

	JSONError(w, l, errors.E(err, errors.Forbidden, "Data not found"))

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Equal(t, `{"msg":"Data not found","status":403}`, w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.HeaderMap.Get("Content-Type"))
}
