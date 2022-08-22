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
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPopulate(t *testing.T) {
	type A struct {
		Field0            int
		Field1            int
		Field2            string
		Field3            *time.Time
		Field4            string
		fieldCanNotSet    string
		FieldTagSkip      string
		FieldTagNameFromB string
	}

	type B struct {
		Field1            int
		Field2            *string
		Field3            time.Time
		field4            string
		Field5            string
		FieldTagSkip      string `populate:"-"`
		FieldByTagNameToA string `populate:"FieldTagNameFromB"`
	}

	fTime := time.Now()
	a := A{
		Field1:       5,
		Field2:       "no-update",
		Field4:       "un-exported-b",
		FieldTagSkip: "skipped",
		Field3:       &fTime,
	}
	b := B{
		Field1:            1,
		Field3:            fTime,
		field4:            "un-exported",
		Field5:            "not in A",
		FieldTagSkip:      "skip-Pfailed",
		FieldByTagNameToA: "from-b-tag-name",
	}
	want := A{
		Field0:            0,
		Field1:            1,
		Field2:            "no-update",
		Field3:            &fTime,
		Field4:            "un-exported-b",
		FieldTagSkip:      "skipped",
		FieldTagNameFromB: "from-b-tag-name",
	}

	// a, b pointers
	err := Populate(&a, &b)
	assert.Nil(t, err)
	assert.Equal(t, want, a)

	// just a pointer
	err = Populate(&a, b)
	assert.Nil(t, err)
	assert.Equal(t, want, a)

	// a pointer, b not
	err = Populate(&a, b)
	assert.Nil(t, err)
	assert.Equal(t, want, a)

	//
	// test fails
	//

	// must pass point for first struct
	a = A{
		Field1: 5,
		Field2: "no-update",
		Field4: "un-exported-b",
	}
	b = B{
		Field1: 1,
		Field3: fTime,
		field4: "un-exported",
		Field5: "not in A",
	}

	err = Populate(a, &b)
	assert.NotNil(t, err)

	err = Populate(nil, b)
	assert.NotNil(t, err)

	err = Populate(a, nil)
	assert.NotNil(t, err)

	notStruct := "not struct"
	err = Populate(&notStruct, "not struct")
	assert.NotNil(t, err)

	err = Populate(&a, "not struct")
	assert.NotNil(t, err)
}
