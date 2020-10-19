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
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoArgs(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("E() did not panic")
		}
	}()
	_ = E()
}

func TestError_Error(t *testing.T) {
	tests := map[string]struct {
		err     error
		match   string
		matched bool
	}{
		"no error":          {&Error{}, "no error", true},
		"skip HTTP message": {E("HTTP message"), "HTTP message", false},
		"just kind":         {E(Unprocessable), ": unprocessable", true},

		"just error wrap":          {E(fmt.Errorf("wrapped error")), ": wrapped error", true},
		"wrap, kind, HTTP message": {E(fmt.Errorf("wrapped error"), Unprocessable, "HTTP message"), ": unprocessable: wrapped error", true},

		// Nested *Error values.
		"nesting-wrap, kind":          {E(E(fmt.Errorf("wrapped error"), Unprocessable), NotFound), ":: errors_test.go", true},
		"nesting-wrap, no-kind, kind": {E(E(fmt.Errorf("wrapped error"), Unprocessable)), ":: errors_test.go", true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := test.err.Error()
			if !assert.Equal(t, test.matched, strings.Contains(got, test.match)) {
				t.Errorf("got: %s; want %s", got, test.match)
			}
		})
	}
}

func TestKind(t *testing.T) {
	tests := map[string]struct {
		err  error
		kind Kind
		want bool
	}{
		// Non-Error errors.
		"nil":        {nil, Unprocessable, false},
		"not *Error": {fmt.Errorf("not an *Error"), Unprocessable, false},

		// Basic comparisons.
		"other":         {E(Other), Unprocessable, false},
		"unprocessable": {E(Unprocessable), Unprocessable, true},
		"internal":      {E(Internal), Unprocessable, false},
		"no kind":       {E(1), Unprocessable, false},

		// Nested *Error values.
		"nesting-unprocessable": {E("Nesting", E(Unprocessable)), Unprocessable, true},
		"nesting-internal":      {E("Nesting", E(Internal)), Unprocessable, false},
		"nesting-no kind":       {E("Nesting", E(1)), Unprocessable, false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := IsKind(test.kind, test.err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestKind_String(t *testing.T) {
	tests := map[string]struct {
		kind Kind
		want string
	}{
		"Other":         {Other, "other error"},
		"Unauthorized":  {Unauthorized, "unauthorized"},
		"Forbidden":     {Forbidden, "forbidden"},
		"NotFound":      {NotFound, "not found"},
		"Conflict":      {Conflict, "conflict"},
		"Unprocessable": {Unprocessable, "unprocessable"},
		"Internal":      {Internal, "internal error"},
		"unknown":       {Kind(999), "unknown error kind"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := test.kind.String()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestToHTTPResponse(t *testing.T) {

	tests := map[string]struct {
		err  *Error
		want string
	}{
		"no message":       {&Error{HTTPMessage: ""}, ""},
		"message":          {&Error{HTTPMessage: "message"}, "message"},
		"chained messages": {&Error{HTTPMessage: "message", Err: &Error{HTTPMessage: "message 2", Err: &Error{HTTPMessage: "message 3"}}}, "message: message 2: message 3"},
		"chained messages, suppress consecutive duplications": {&Error{HTTPMessage: "message", Err: &Error{HTTPMessage: "message"}}, "message"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := ToHTTPResponse(test.err)
			assert.Equal(t, test.want, got)
		})
	}
}
