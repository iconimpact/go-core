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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMailosaur_MessageByID(t *testing.T) {
	if !envVars {
		t.Skip(noEnvVarsMsg)
	}

	message, err := ms.MessageByID(testMessages.Items[0].ID)
	assert.Nil(t, err)

	assert.Equal(t, "Go mailosaur", message.Subject)
	assert.Equal(t, "Go mailosaur!", message.Text.Body)

	// test fails
	_, err = ms.MessageByID("no-message-id")
	assert.NotNil(t, err)
}

func TestMailosaur_MessageDelete(t *testing.T) {
	if !envVars {
		t.Skip(noEnvVarsMsg)
	}

	err := ms.MessageDelete(testMessages.Items[0].ID)
	assert.Nil(t, err)

	// validate message deleted
	_, err = ms.MessageByID(testMessages.Items[0].ID)
	assert.Contains(t, err.Error(), "404")

	// test fails
	err = ms.MessageDelete("no-message-id")
	assert.NotNil(t, err)
}

func TestMailosaur_MessagesList(t *testing.T) {
	if !envVars {
		t.Skip(noEnvVarsMsg)
	}

	messages, err := ms.MessagesList(serverID)
	assert.Nil(t, err)

	if assert.True(t, len(messages.Items) > 0) {
		assert.Equal(t, "Go mailosaur", messages.Items[0].Subject)
		assert.Equal(t, "", messages.Items[0].Text.Body)
	}

	// test fails
	_, err = ms.MessagesList("no-server-id")
	assert.NotNil(t, err)
}

func TestMailosaur_MessagesSearch(t *testing.T) {
	if !envVars {
		t.Skip(noEnvVarsMsg)
	}

	searchCriteria := SearchCriteria{
		SentTo:  sendTo,
		Subject: "Go mailosaur",
		Match:   "ALL",
	}

	messages, err := ms.MessagesSearch(serverID, searchCriteria)
	assert.Nil(t, err)

	if assert.True(t, len(messages.Items) > 0) {
		assert.Equal(t, "Go mailosaur", messages.Items[0].Subject)
		assert.Equal(t, "", messages.Items[0].Text.Body)
	}

	// test fails
	_, err = ms.MessagesSearch("no-server-id", searchCriteria)
	assert.NotNil(t, err)

	_, err = ms.MessagesSearch(serverID, SearchCriteria{})
	assert.NotNil(t, err)
}

func TestMailosaur_MessagesDelete(t *testing.T) {
	if !envVars {
		t.Skip(noEnvVarsMsg)
	}

	err := ms.MessagesDelete(serverID)
	assert.Nil(t, err)

	messages, err := ms.MessagesList(serverID)
	assert.Nil(t, err)

	assert.True(t, len(messages.Items) == 0)

	// test fails
	err = ms.MessagesDelete("no-server-id")
	assert.NotNil(t, err)
}
