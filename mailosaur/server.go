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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Server represent a Mailosaur test server
// https://mailosaur.com/docs/api/reference/servers/
type Server struct {
	ID              string   `json:"id"`
	Password        string   `json:"password"`
	Name            string   `json:"name"`
	Users           []string `json:"users"`
	Messages        int      `json:"messages"`
	ForwardingRules []string `json:"forwardingRules"`
	Retention       int      `json:"retention"`
}

// ServerByID retrieves the detail for a server by ID.
// https://mailosaur.com/docs/api/reference/servers/#retrieve-a-server
// GET https://mailosaur.com/api/servers/:id
func (m *Mailosaur) ServerByID(ID string) (*Server, error) {
	url := fmt.Sprintf("%s/servers/%s", baseURL, url.QueryEscape(ID))

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
		err = fmt.Errorf("mailosaur.ServerByID failed with status: %d, response: %s", resp.StatusCode, body)
		return nil, err
	}

	server := &Server{}
	err = json.Unmarshal(body, server)
	if err != nil {
		return nil, err
	}

	return server, nil
}
