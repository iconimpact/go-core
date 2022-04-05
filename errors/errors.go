/*
   Copyright 2020 iconmobile GmbH

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package errors

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path"
	"runtime"
)

// Kind defines the error type this is, mostly for use by routers
// that must set different response status depending on the error.
type Kind int

// Kinds types.
const (
	Other         Kind = iota // Unclassified error
	BadRequest                // Bad Request (400)
	Unauthorized              // Unauthorized (401)
	Forbidden                 // Forbidden (403)
	NotFound                  // Not found (404)
	Conflict                  // Conflict (409)
	Gone                      // Gone (410)
	Unprocessable             // Unprocessable, invalid request data (422)
	Internal                  // Internal server error (500)
	BadGateway                // Bad gateway (502)
)

// Separator defines the string used to separate nested errors.
var Separator = " »» "

// String transforms Kind type to text.
func (k Kind) String() string {
	switch k {
	case Other:
		return "other error"
	case BadRequest:
		return "bad request"
	case Unauthorized:
		return "unauthorized"
	case Forbidden:
		return "forbidden"
	case NotFound:
		return "not found"
	case Conflict:
		return "conflict"
	case Gone:
		return "gone"
	case Unprocessable:
		return "unprocessable"
	case Internal:
		return "internal error"
	case BadGateway:
		return "bad gateway"
	}
	return "unknown error kind"
}

// Error defines a standard application error.
type Error struct {
	// application specific fields.
	HTTPMessage string

	// logical operation and nested error.
	Kind Kind
	Err  error

	// stack information.
	stack string
}

func (e *Error) isZero() bool {
	return e.HTTPMessage == "" && e.Kind == 0 && e.Err == nil
}

// pad appends str to the buffer if the buffer already has some data.
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}

func (e *Error) Error() string {
	b := new(bytes.Buffer)
	if e.stack != "" {
		pad(b, ": ")
		b.WriteString(e.stack)
	}
	if e.Kind != 0 {
		pad(b, ": ")
		b.WriteString(e.Kind.String())
	}
	if e.Err != nil {
		// if we are cascading non-empty Error errors custom pad.
		if prevErr, ok := e.Err.(*Error); ok {
			if !prevErr.isZero() {
				pad(b, Separator)
				b.WriteString(e.Err.Error())
			}
		} else {
			pad(b, ": ")
			b.WriteString(e.Err.Error())
		}
	}
	if b.Len() == 0 {
		return "no error"
	}
	return b.String()
}

// E builds an error value from its arguments.
// There must be at least one argument or E panics.
// The type of each argument determines its meaning.
// If more than one argument of a given type is presented,
// only the last one is recorded.
//
// The types are:
//	string
//		The HTTP message for the API user.
//	errors.Kind
//		The class of error, such as permission failure.
//	error
//		The underlying error that triggered this one.
//
// If the error is printed, only those items that have been
// set to non-zero values will appear in the result.
//
// If Kind is not specified or Other, we set it to the Kind of
// the underlying error.
//
func E(args ...interface{}) error {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}

	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			e.HTTPMessage = arg
		case Kind:
			e.Kind = arg
		case *Error:
			copy := *arg
			e.Err = &copy
		case error:
			e.Err = arg
		default:
			_, file, line, _ := runtime.Caller(1)
			log.Printf("errors.E: bad call from %s:%d: %v", file, line, args)
			return fmt.Errorf("unknown type %T, value %v in error call", arg, arg)
		}
	}

	// record stack
	_, file, line, ok := runtime.Caller(1)
	if ok {
		e.stack = fmt.Sprintf("%s:%d", getFileName(file), line)
	}

	prev, ok := e.Err.(*Error)
	if !ok {
		return e
	}

	// The previous error was also one of ours. Suppress duplications
	// so the message won't contain the same kind, HTTP message
	// twice.
	if prev.HTTPMessage == e.HTTPMessage {
		prev.HTTPMessage = ""
	}

	// If this error has Kind unset or Other, pull up the inner one.
	if e.Kind == Other {
		e.Kind = prev.Kind
		prev.Kind = Other
	}

	return e
}

func getFileName(file string) string {
	_, fileName := path.Split(file)
	return fileName
}

// ToHTTPStatus converts an error to an HTTP status code.
func ToHTTPStatus(e *Error) int {
	if e == nil {
		return http.StatusInternalServerError
	}

	var status int
	switch e.Kind {
	case BadRequest:
		status = http.StatusBadRequest
	case Unauthorized:
		status = http.StatusUnauthorized
	case Forbidden:
		status = http.StatusForbidden
	case NotFound:
		status = http.StatusNotFound
	case Conflict:
		status = http.StatusConflict
	case Gone:
		status = http.StatusGone
	case Unprocessable:
		status = http.StatusUnprocessableEntity
	case Internal:
		status = http.StatusInternalServerError
	case BadGateway:
		status = http.StatusBadGateway
	default:
		status = http.StatusInternalServerError
	}

	return status
}

// ToHTTPResponse creates a string to be used for HTTP response
// by chaining the underlying application errors HTTPMessage.
func ToHTTPResponse(e *Error) string {
	if e == nil {
		return ""
	}

	b := new(bytes.Buffer)

	if e.HTTPMessage != "" {
		pad(b, ": ")
		b.WriteString(e.HTTPMessage)
	}

	prev, ok := e.Err.(*Error)
	if !ok {
		return b.String()
	}

	// suppress consecutive duplications
	if prev.HTTPMessage == e.HTTPMessage {
		prev.HTTPMessage = ""
	}

	// add pad
	if prev.HTTPMessage != "" {
		pad(b, ": ")
	}

	b.WriteString(ToHTTPResponse(prev))

	return b.String()
}

// IsKind reports whether err is an *Error of the given Kind.
// If err is nil then Is returns false.
func IsKind(kind Kind, err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.Kind != Other {
		return e.Kind == kind
	}
	if e.Err != nil {
		return IsKind(kind, e.Err)
	}
	return false
}

// Is provides compatibility for Go 1.13 error chains.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return (e.Kind == t.Kind || t.Kind == 0) &&
		(Is(e.Err, t.Err) || t.Err == nil)
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *Error) Unwrap() error { return e.Err }
