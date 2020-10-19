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
	"fmt"
	"reflect"
)

const populateTagName = "populate"

// Populate sets fields from 'b' into 'a' struct by matching name and type,
// use pointers in 'b' fields to pass empty values.
func Populate(a, b interface{}) error {
	// a must be struct
	va := reflect.ValueOf(a)
	if va.Kind() != reflect.Ptr {
		return fmt.Errorf("must pass a pointer, not a value as a struct")
	}
	va = va.Elem()

	if va.Kind() != reflect.Struct {
		return fmt.Errorf("a must be struct")
	}

	// b must be struct
	vb := reflect.ValueOf(b)
	if vb.Kind() == reflect.Ptr {
		vb = vb.Elem()
	}
	if vb.Kind() != reflect.Struct {
		return fmt.Errorf("b must be struct")
	}

	t := vb.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Name
		value := vb.FieldByName(name)

		// skip un-exported
		if !isExported(field) {
			continue
		}

		// check tag
		tagFieldName, _ := parseTag(field.Tag.Get(populateTagName))
		if tagFieldName != "" {
			name = tagFieldName
		}

		// "-" tag ignores the field
		if name == "-" {
			continue
		}

		// check dest field exist
		// and can be set
		aField, ok := hasField(a, name)
		if !ok {
			continue
		}
		if !aField.CanSet() {
			continue
		}

		// skip nil
		// to accept empty value
		// field must be pointer
		if isZero(value) {
			continue
		}

		// switch pointer type
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}

		// update dest field
		// if types match
		if aField.Type() == value.Type() {
			aField.Set(value)
		}
	}

	return nil
}

func hasField(s interface{}, name string) (reflect.Value, bool) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	field := v.FieldByName(name)
	if !field.IsValid() {
		return field, false
	}

	return field, true
}

func isZero(value reflect.Value) bool {
	zero := reflect.Zero(value.Type()).Interface()
	current := value.Interface()

	return reflect.DeepEqual(current, zero)
}

func isExported(field reflect.StructField) bool {
	return field.PkgPath == ""
}
