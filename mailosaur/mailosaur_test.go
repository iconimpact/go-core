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
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"gopkg.in/gomail.v2"
)

var ms *Mailosaur
var serverID string
var sendTo string
var testMessages *SearchResponse

func TestMain(m *testing.M) {
	// setup before tests
	apiKey, ok := os.LookupEnv("MAILOSAUR_API_KEY")
	if !ok {
		log.Println("env var MAILOSAUR_API_KEY required")
		os.Exit(1)
	}

	serverID, ok = os.LookupEnv("MAILOSAUR_SERVER")
	if !ok {
		log.Println("env var MAILOSAUR_SERVER required")
		os.Exit(1)
	}

	// init mailosaur client
	ms = New(apiKey)

	// create an email
	server, err := ms.ServerByID(serverID)
	if err != nil {
		log.Printf("server ID info err: %v\n", err)
		os.Exit(1)
	}

	username := fmt.Sprintf("%s@mailosaur.io", server.ID)
	password := server.Password
	sendTo = fmt.Sprintf("gotest.%s@mailosaur.io", server.ID)

	err = sendMail(username, password, sendTo)
	if err != nil {
		log.Printf("sendMail err: %v\n", err)
		os.Exit(1)
	}
	err = sendMail(username, password, sendTo)
	if err != nil {
		log.Printf("sendMail err: %v\n", err)
		os.Exit(1)
	}
	// wait for the emails to be received
	time.Sleep(7 * time.Second)

	// get created mails
	searchCriteria := SearchCriteria{
		SentTo:  sendTo,
		Subject: "Go mailosaur",
		Match:   "ALL",
	}
	testMessages, err = ms.MessagesSearch(serverID, searchCriteria)
	if err != nil || len(testMessages.Items) == 0 {
		log.Printf("ms.MessagesSearch err: %v\n", err)
		os.Exit(1)
	}

	// run tests
	code := m.Run()

	// cleanup after tests

	os.Exit(code)
}

func sendMail(username, password, sendTo string) error {
	host := "mailosaur.io"
	port := 587

	d := gomail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", sendTo)
	m.SetHeader("Subject", "Go mailosaur")
	m.SetBody("text/plain", "Go mailosaur!")

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
