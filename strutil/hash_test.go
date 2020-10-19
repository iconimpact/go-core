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

/*
About Benchmarking:

call with (1 second default time):
go test -bench=. -cpuprofile profile_cpu.out
go test -bench=. -cpuprofile profile_cpu.out -test.benchtime 10s

Learn more at https://scene-si.org/2017/06/06/benchmarking-go-programs/

More easy to read output with prettybench:

go get github.com/cespare/prettybench
go test -bench=. | prettybench
*/

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	text := "Hi, I want to be hashed"
	description := "a test hash used in secure_test.go"
	description2 := description + "... not"
	hash := "b977e11174549a4ed5b73f88ae80d2d3ae155" +
		"4f5da87fb6bfd2fa84c6ca53fee"
	assert.Equal(t, hash, Hash(text, description))
	assert.NotEqual(t, hash, Hash(text, description2))
	assert.NotEqual(t, hash, Hash(text+"x", description))
	assert.Equal(t, 64, len(Hash(text, description)))
	assert.Equal(t, 64, len(Hash(text, description2)))
}

//////////////////////
// Benchmarks
//////////////////////

func BenchmarkHash(b *testing.B) {
	b.ReportAllocs()
	password := "My dark, little secret"

	for i := 0; i < b.N; i++ {
		res := Hash(password, "test password")
		if res == "" {
			b.Error("Unexpected result: " + res)
		}
	}
}

//////////////////////
// TEST HELPERS
//////////////////////

// should be called with a defer command inside logic
func trackTime(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
