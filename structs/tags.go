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
	"strings"
)

// tagOptions tag options slice
type tagOptions []string

// Has returns true if opt in tagOptions.
func (t tagOptions) Has(opt string) bool {
	for _, tagOpt := range t {
		if tagOpt == opt {
			return true
		}
	}

	return false
}

// parseTag returns a struct field's tag name and list of options.
// A tag is in the form of: "name,option1,option2".
// The name can be neglectected.
func parseTag(tag string) (string, tagOptions) {
	parts := strings.Split(tag, ",")
	return parts[0], parts[1:]
}
