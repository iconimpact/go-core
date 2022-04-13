package respond

import (
	"net/http"

	"github.com/iconmobile-dev/go-core/errors"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

const jsonContentType = "application/json; charset=utf-8"

var json = jsoniter.ConfigFastest

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json" and
// X-Content-Type-Options as "nosniff".
// Logs the status and v if l is not nil.
func JSON(w http.ResponseWriter, l *zap.Logger, status int, v interface{}) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		panic("respond: " + err.Error())
	}

	if l != nil {
		l.Info("respond: ", zap.Int("status", status), zap.String("body", string(jsonBytes)))
	}

	w.Header().Set("Content-Type", jsonContentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	_, err = w.Write(jsonBytes)
	if err != nil {
		panic("respond: " + err.Error())
	}
}

// JSONError returns an HTTP response as JSON message with status code
// base on app err Kind, Msg from app err HTTPMessage.
// Logs the error if l is not nil.
func JSONError(w http.ResponseWriter, l *zap.Logger, err error) {
	var status int
	var errMsg string

	if l != nil {
		l.Error("respond: ", zap.Error(err))
	}

	// set custom app err Message
	appErr, ok := err.(*errors.Error)
	if !ok {
		status = http.StatusInternalServerError
		errMsg = "Internal Server Error"
	} else {
		status = errors.ToHTTPStatus(appErr)
		errMsg = errors.ToHTTPResponse(appErr)
	}

	JSON(w, nil, status, errMsg)
}
