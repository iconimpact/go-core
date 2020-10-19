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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Message represent an email or SMS received by Mailosaur
// and contain all the data you might need to perform any number of manual or automated tests.
// https://mailosaur.com/docs/api/reference/messages/#the-message-object
type Message struct {
	ID     string `json:"id"`
	Server string `json:"server"`
	Rcpt   []struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"rcpt"`
	From []struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"from"`
	To []struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"to"`
	Cc []struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"cc"`
	Bcc []struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"bcc"`
	Received time.Time `json:"received"`
	Subject  string    `json:"subject"`
	HTML     struct {
		Links  []string `json:"links"`
		Images []string `json:"images"`
		Body   string   `json:"body"`
	} `json:"html"`
	Text struct {
		Links  []string `json:"links"`
		Images []string `json:"images"`
		Body   string   `json:"body"`
	} `json:"text"`
	Attachments []struct {
		ID          string `json:"id"`
		ContentType string `json:"contentType"`
		FileName    string `json:"fileName"`
		ContentID   string `json:"contentId"`
		Length      string `json:"length"`
		URL         string `json:"url"`
	} `json:"attachments"`
	Metadata struct {
		Headers []struct {
			Field string `json:"field"`
			Value string `json:"value"`
		} `json:"headers"`
	} `json:"metadata"`
	HateosLinks []struct {
		Href   string `json:"href"`
		Method string `json:"method"`
		Rel    string `json:"rel"`
	} `json:"hateosLinks"`
	Read bool `json:"read"`
}

// UnmarshalJSON is a custom UnmarshalJSON to deal with search respons
// of attachments as 0
func (m *Message) UnmarshalJSON(data []byte) error {
	type Alias Message

	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	va, ok := v["attachments"]
	if !ok || va == float64(0) {
		aux := &struct {
			Attachments interface{} `json:"attachments"`
			*Alias
		}{
			Alias: (*Alias)(m),
		}

		if err := json.Unmarshal(data, &aux); err != nil {
			return err
		}

		return nil
	}

	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	return nil
}

// MessageByID retrieves the detail for a single email message by ID.
// https://mailosaur.com/docs/api/reference/messages/#retrieve-a-message
// GET https://mailosaur.com/api/messages/:id
func (m *Mailosaur) MessageByID(ID string) (*Message, error) {
	url := fmt.Sprintf("%s/messages/%s", baseURL, url.QueryEscape(ID))

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "en")
	req.SetBasicAuth(m.apiKey, "")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err = fmt.Errorf("mailosaur.MessageByID failed with status: %d, response: %s", resp.StatusCode, body)
		return nil, err
	}

	message := &Message{}
	err = json.Unmarshal(body, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// MessageDelete permanently deletes a message and any attachments related to the message.
// https://mailosaur.com/docs/api/reference/messages/#delete-a-message
// DELETE https://mailosaur.com/api/messages/:id
func (m *Mailosaur) MessageDelete(ID string) error {
	url := fmt.Sprintf("%s/messages/%s", baseURL, url.QueryEscape(ID))

	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "en")
	req.SetBasicAuth(m.apiKey, "")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = fmt.Errorf("mailosaur.MessageDelete failed with status: %d, response: %s", resp.StatusCode, body)
		return err
	}

	return nil
}

// ListResponse response from MessagesList.
type ListResponse struct {
	Items []Message `json:"items"`
}

// MessagesList returns a list of your messages in summary form.
// The summaries are returned sorted by received date, with the most recently-received messages appearing first.
// https://mailosaur.com/docs/api/reference/messages/#list-all-messages
// GET https://mailosaur.com/api/messages?server=:server
func (m *Mailosaur) MessagesList(server string, opt ...CallOption) (*ListResponse, error) {
	url := fmt.Sprintf("%s/messages?server=%s", baseURL, url.QueryEscape(server))

	u, err := addOptions(url, opt...)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "en")
	req.SetBasicAuth(m.apiKey, "")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err = fmt.Errorf("mailosaur.MessagesList failed with status: %d, response: %s", resp.StatusCode, body)
		return nil, err
	}

	messages := &ListResponse{}
	err = json.Unmarshal(body, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// MessagesDelete permanently deletes all messages held by the specified server.
// Also deletes any attachments related to each message.
// https://mailosaur.com/docs/api/reference/messages/#delete-all-messages
// DELETE https://mailosaur.com/api/messages?server=:server
func (m *Mailosaur) MessagesDelete(server string) error {
	url := fmt.Sprintf("%s/messages?server=%s", baseURL, url.QueryEscape(server))

	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "en")
	req.SetBasicAuth(m.apiKey, "")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = fmt.Errorf("mailosaur.MessagesDelete failed with status: %d, response: %s", resp.StatusCode, body)
		return err
	}

	return nil
}

// SearchResponse response from MessagesSearch.
type SearchResponse struct {
	Items []Message `json:"items"`
}

// MessagesSearch returns a list of message summaries matching the specified search criteria, in summary form.
// The summaries are returned sorted by received date, with the most recently-received messages appearing first.
// https://mailosaur.com/docs/api/reference/messages/#search-for-messages
// POST https://mailosaur.com/api/messages/search?server=:server
func (m *Mailosaur) MessagesSearch(server string, searchCriteria SearchCriteria, opt ...CallOption) (*SearchResponse, error) {
	url := fmt.Sprintf("%s/messages/search?server=%s", baseURL, url.QueryEscape(server))

	u, err := addOptions(url, opt...)
	if err != nil {
		return nil, err
	}

	jsonStr, err := json.Marshal(searchCriteria)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", u, bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "en")
	req.SetBasicAuth(m.apiKey, "")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err = fmt.Errorf("mailosaur.MessagesSearch failed with status: %d, response: %s", resp.StatusCode, body)
		return nil, err
	}

	searchResponse := &SearchResponse{}
	err = json.Unmarshal(body, searchResponse)
	if err != nil {
		return nil, err
	}

	return searchResponse, nil
}
