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
	"regexp"
	"strings"
)

const sanitizeTagName = "sanitize"

// Sanitize removes all leading and trailing white spaces
// from struct exported string fields
func Sanitize(s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("must pass a pointer, not a value")
	}
	v = v.Elem()

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("must pass a struct")
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Name

		// check tag
		tagFieldName, _ := parseTag(field.Tag.Get(sanitizeTagName))

		// "-" tag ignores the field
		if tagFieldName == "-" {
			continue
		}

		sv := v.FieldByName(name)
		if sv.Kind() == reflect.Ptr {
			sv = sv.Elem()
		}

		// check if field can be reset
		if !sv.CanSet() {
			continue
		}

		// only string fields
		if sv.Kind() != reflect.String {
			continue
		}

		svt := strings.TrimSpace(sv.Interface().(string))

		// "email" tag cleans, validates email
		if tagFieldName == "email" {
			var err error
			svt, err = CleanEmail(svt)
			if err != nil {
				return err
			}
		}

		sv.SetString(svt)
	}
	return nil
}

// IsEmail checks if an email is RFC 2821, 2822 compliant
// taken from http://regexlib.com/REDetails.aspx?regexp_id=2558
// returns true if it is a valid email
func IsEmail(email string) bool {
	r := "^((([!#$%&'*+\\-/=?^_`{|}~\\w])|([!#$%&'*+\\-/=?^_`{|}~\\w]"
	r += "[!#$%&'*+\\-/=?^_`{|}~\\.\\w]{0,}[!#$%&'*+\\-/=?^_`{|}~\\w]))"
	r += "[@]\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*)$"
	re := regexp.MustCompile(r)
	return re.MatchString(email)
}

// CleanEmail regexp cleans and validates string is an email address
// and all @googlemail.com addresses are
// normalized to @gmail.com.
func CleanEmail(str string) (string, error) {
	str = strings.ToLower(str)
	str = strings.TrimSpace(str)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "@googlemail.", "@gmail.", -1)

	reg, err := regexp.Compile("[^a-zA-Z0-9_@.+-]+")
	if err != nil {
		return "", fmt.Errorf("regexp email compile: %w", err)
	}
	str = reg.ReplaceAllString(str, "")

	if !IsEmail(str) {
		return "", fmt.Errorf("%s is not an email", str)
	}

	return str, nil
}
