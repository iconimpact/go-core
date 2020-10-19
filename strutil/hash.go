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
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

// Hash generates a hash of data using HMAC-SHA-512/256 and adds
// the description text as HMAC "key" to get different hashes
// for different purposes.
// This function is NOT SECURE ENOUGH for securely
// hashing normal user-entered passwords, use HashPassword instead!
func Hash(text string, description string) string {
	// hash some rounds
	// it is still incredibly fast (0.7 micro seconds)
	rounds := 3
	hash := text
	for i := 1; i <= rounds; i++ {
		h := hmac.New(sha512.New512_256, []byte(description))
		h.Write([]byte(hash))
		hash = hex.EncodeToString(h.Sum(nil))
	}
	return hash
}

func sha256Hash(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}
