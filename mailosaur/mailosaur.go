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

package mailosaur

import (
	"fmt"
	"net/url"
)

const baseURL = "https://mailosaur.com/api"

// Mailosaur client to work with Mailosaur API.
// https://mailosaur.com/docs/api/
type Mailosaur struct {
	apiKey string
}

// New init the mailosaur client
func New(apiKey string) *Mailosaur {
	return &Mailosaur{
		apiKey: apiKey,
	}
}

// CallOption is an optional argument to an API call:
// page, itemsPerPage and etc
type CallOption interface {
	Get() (key, value string)
}

// OptServer is an optional argument to an API call
type OptServer string

// Get return key/value for make query
func (o OptServer) Get() (string, string) {
	return "server", fmt.Sprint(o)
}

// OptPage is an optional argument to an API call
type OptPage string

// Get return key/value for make query
func (o OptPage) Get() (string, string) {
	return "page", fmt.Sprint(o)
}

// OptItemsPerPage is an optional argument to an API call
type OptItemsPerPage string

// Get return key/value for make query
func (o OptItemsPerPage) Get() (string, string) {
	return "itemsPerPage", fmt.Sprint(o)
}

// OptReceivedAfter is an optional argument to an API call
type OptReceivedAfter string

// Get return key/value for make query
func (o OptReceivedAfter) Get() (string, string) {
	return "receivedAfter", fmt.Sprint(o)
}

// OptSentTo is an optional argument to an API call
type OptSentTo string

// Get return key/value for make query
func (o OptSentTo) Get() (string, string) {
	return "sentTo", fmt.Sprint(o)
}

// OptSubject is an optional argument to an API call
type OptSubject string

// Get return key/value for make query
func (o OptSubject) Get() (string, string) {
	return "subject", fmt.Sprint(o)
}

// OptBody is an optional argument to an API call
type OptBody string

// Get return key/value for make query
func (o OptBody) Get() (string, string) {
	return "body", fmt.Sprint(o)
}

// OptMatch is an optional argument to an API call
type OptMatch string

// Get return key/value for make query
func (o OptMatch) Get() (string, string) {
	return "match", fmt.Sprint(o)
}

func addOptions(s string, opts ...CallOption) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs := u.Query()
	for _, o := range opts {
		qs.Set(o.Get())
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// SearchCriteria to search for a message,
// match if set to ALL (default), then only results that match all specified criteria will be returned.
// If set to ANY, results that match any of the specified criteria will be returned.
type SearchCriteria struct {
	SentTo  string `json:"sentTo"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Match   string `json:"match"`
}
