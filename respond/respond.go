package respond

import (
	"bytes"
	"io"
	"net/http"
	"reflect"

	"github.com/iconmobile-dev/go-core/errors"
	jsoniter "github.com/json-iterator/go"
)

const jsonContentType = "application/json; charset=utf-8"

var json = jsoniter.ConfigFastest

// JSONMsg default struct for message response body.
type JSONMsg struct {
	Data interface{} `json:"data" swaggertype:"object"`
	Msg  string      `json:"msg"`
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json" and
// X-Content-Type-Options as "nosniff".
// If v is a string then it is sent as "Msg" property value.
func JSON(w http.ResponseWriter, r *http.Request, status int, v interface{}) {
	var body []byte
	var data JSONMsg

	iv := reflect.ValueOf(v)
	if iv.Kind() == reflect.Ptr {
		iv = iv.Elem()
	}

	switch iv.Kind() {
	case reflect.String:
		data.Msg = v.(string)
	default:
		data.Data = v
	}

	body, err := json.Marshal(&data)
	if err != nil {
		panic("respond: " + err.Error())
	}

	w.Header().Set("Content-Type", jsonContentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	_, err = io.Copy(w, bytes.NewReader(body))
	if err != nil {
		panic("respond: " + err.Error())
	}
}

// JSONError returns an HTTP response as JSON message with status code
// base on app err Kind, Msg from app err HTTPMessage.
// Logs the error by default.
func JSONError(w http.ResponseWriter, r *http.Request, err error) {
	var status int
	var errMsg string

	// log the error by default
	logger.Error("respond: " + err.Error())

	// set custom app err Message
	appErr, ok := err.(*errors.Error)
	if !ok {
		status = http.StatusInternalServerError
		errMsg = "Internal Server Error"
	} else {
		switch appErr.Kind {
		case errors.BadRequest:
			status = http.StatusBadRequest
		case errors.Unauthorized:
			status = http.StatusUnauthorized
		case errors.Forbidden:
			status = http.StatusForbidden
		case errors.NotFound:
			status = http.StatusNotFound
		case errors.Conflict:
			status = http.StatusConflict
		case errors.Gone:
			status = http.StatusGone
		case errors.Unprocessable:
			status = http.StatusUnprocessableEntity
		case errors.Internal:
			status = http.StatusInternalServerError
		case errors.BadGateway:
			status = http.StatusBadGateway
		default:
			status = http.StatusInternalServerError
		}

		errMsg = errors.ToHTTPResponse(appErr)
	}

	JSON(w, r, status, errMsg)
}
