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

import "testing"

func TestParseTag_Name(t *testing.T) {
	tags := []struct {
		tag string
		has bool
	}{
		{"", false},
		{"name", true},
		{"name,opt", true},
		{"name , opt, opt2", false}, // has a single whitespace
		{", opt, opt2", false},
	}

	for _, tag := range tags {
		name, _ := parseTag(tag.tag)

		if (name != "name") && tag.has {
			t.Errorf("Parse tag should return name: %#v", tag)
		}
	}
}

func TestParseTag_Opts(t *testing.T) {
	tags := []struct {
		opts string
		has  bool
	}{
		{"name", false},
		{"name,opt", true},
		{"name , opt, opt2", false}, // has a single whitespace
		{",opt, opt2", true},
		{", opt3, opt4", false},
	}

	// search for "opt"
	for _, tag := range tags {
		_, opts := parseTag(tag.opts)

		if opts.Has("opt") != tag.has {
			t.Errorf("Tag opts should have opt: %#v", tag)
		}
	}
}
