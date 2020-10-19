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

package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	type S struct {
		SToTrim     string
		PointToS    *string
		Integer     int
		sUnExported string
		TagSkip     string `sanitize:"-"`
		TagEmail    string `sanitize:"email"`
	}

	pointerString := " no-trim "
	s := &S{
		SToTrim:     " trim ",
		PointToS:    &pointerString,
		Integer:     10,
		sUnExported: " un-exported ",
		TagSkip:     " skipped ",
		TagEmail:    " email.string+me@googlemail.com",
	}

	want := &S{
		SToTrim:     "trim",
		PointToS:    &pointerString,
		Integer:     10,
		sUnExported: " un-exported ",
		TagSkip:     " skipped ",
		TagEmail:    "email.string+me@gmail.com",
	}

	err := Sanitize(s)
	assert.Nil(t, err)
	assert.Equal(t, want, s)

	//
	// test fails
	//
	err = Sanitize(*s)
	assert.NotNil(t, err)

	notStruct := "not struct"
	err = Sanitize(&notStruct)
	assert.NotNil(t, err)

	s.TagEmail = "not-an-email"
	err = Sanitize(s)
	assert.NotNil(t, err)
}
