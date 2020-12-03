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

package strutil

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	assert.Equal(t, 32, len(Random(32)))
	assert.NotEqual(t, Random(32), Random(32))
}

func TestRandomSecure(t *testing.T) {
	// it is random
	assert.NotEqual(t, RandomSecure(5, ""), RandomSecure(5, ""))

	// length is correct
	assert.Equal(t, 5, len(RandomSecure(5, "")))

	// numbers are returned
	val, err := strconv.Atoi(RandomSecure(5, "number"))
	assert.Nil(t, err)
	assert.True(t, val > 0)

	// alphas (A-Z a-z) are returned
	val, err = strconv.Atoi(RandomSecure(32, "alpha"))
	assert.Error(t, err)

	assert.Equal(t, 32, len(RandomSecure(32, "alpha")))
	assert.Equal(t, 65, len(RandomSecure(65, "alpha")))

	// pin (0-9 A-Z without 0 letter) are returned
	pin := RandomSecure(300, "pin")
	// should not contain 0 letter
	assert.NotContains(t, pin, "O")
	// should not contain lower letter
	assert.NotContains(t, pin, "a")
	assert.NotContains(t, pin, "d")
}
