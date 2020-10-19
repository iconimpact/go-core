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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	text := " My dark, little Secret 🔐 "
	password := "WbuJjNmPPzqi2ikTPeGFXbTqbad9MuP9"
	longPW := password + "x"
	wrongPW := "8JtsDUxzjp372XxNKypHB7zQbnjUBBoG"

	// password is not 32 chars
	encr, err := Encrypt(text, longPW)
	assert.NotNil(t, err)
	assert.Empty(t, encr)

	// decryption fails due to wrong password
	encr, err = Encrypt(text, password)
	assert.Nil(t, err)
	decr, err := Decrypt(encr, wrongPW)
	assert.NotNil(t, err)
	assert.Empty(t, decr)

	// encryption & decryption works
	encr, err = Encrypt(text, password)
	assert.Nil(t, err)
	decr, err = Decrypt(encr, password)
	assert.Equal(t, text, decr)
	assert.Nil(t, err)
}

//////////////////////
// Benchmarks
//////////////////////

func BenchmarkEncryptionDecryption(b *testing.B) {
	b.ReportAllocs()
	text := " My dark, little Secret 🔐 "
	password := "8JtsDUxzjp372XxNKypHB7zQbnjUBBoG"

	for i := 0; i < b.N; i++ {
		encr, err1 := Encrypt(text, password)
		decr, err2 := Decrypt(encr, password)
		if decr != text || err1 != nil || err2 != nil {
			b.Error("Unexpected results")
		}
	}
}
