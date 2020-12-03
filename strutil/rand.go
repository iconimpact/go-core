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
	crand "crypto/rand"
	"fmt"
	mrand "math/rand"
	"time"
)

// Random generates a NOT SECURE! random string of defined length
// The result is NOT SECURE but fast because it uses the time as seed
// taken from https://goo.gl/9GBmNN
func Random(n int) string {
	var src = mrand.NewSource(time.Now().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// RandomSecure returns a SECURELY generated random string
// randType can be:
//	- "alpha" for range A-Z a-z
//	- "number" for range 0-9, returns string number with leading 0s
//	- "pin" for range 0-9 A-Z without O letter
//	- any other value for A-Z a-z 0-9
// It will panic in the super rare case of an issue
// to avoid any cascading security issues
// orients on https://goo.gl/kK987i and https://goo.gl/NRrS7y
func RandomSecure(strSize int, randType string) string {
	dict := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	if randType == "alpha" {
		dict = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	} else if randType == "number" {
		dict = "0123456789"
	} else if randType == "pin" {
		dict = "0123456789ABCDEFGHIJKLMNPQRSTUVWXYZ"
	}

	var b = make([]byte, strSize)
	_, err := crand.Read(b)
	if err != nil {
		msg := "crypto/rand is unavailable: strutil.RandomSecure() "
		msg += "failed with %#v"
		panic(fmt.Sprintf(msg, err))
	}

	// convert the bytes into the appropriate set of chars
	for k, v := range b {
		b[k] = dict[v%byte(len(dict))]
	}

	rndStr1 := string(b)

	if randType == "number" || strSize > 64 {
		return rndStr1
	}

	// if alphanumeric then make it really, really random
	// by taking the current time as seed for a hash
	timeStr := fmt.Sprintf("%d", time.Now().UnixNano())
	rndStr2 := timeStr + rndStr1
	hash := sha256Hash(rndStr2) // 64 chars

	// take the hash from the end
	secret := hash[len(hash)-strSize:]
	return secret
}
