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

package logger

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShouldLogSuccessTxt(t *testing.T) {
	output := &bytes.Buffer{}

	log := Logger{
		MinLevel: "info",
		Output:   output,
	}

	log.Error("Hello there.")

	wanted := "ERROR Hello there."
	assert.Contains(t, output.String(), wanted)
}

func TestNotLogDebugWhenInfoTxt(t *testing.T) {
	output := &bytes.Buffer{}

	log := Logger{
		MinLevel: "info",
		Output:   output,
	}

	log.Debug("Hello there.")

	wanted := ""
	assert.Equal(t, wanted, output.String())
}

func TestDebugJSON(t *testing.T) {
	output := &bytes.Buffer{}

	log := Logger{
		MinLevel: "verbose",
		UseJSON:  true,
		Output:   output,
	}

	log.Debug("Hello there.")

	wanted := `"message":"Hello there.","level_int":1`
	assert.Contains(t, output.String(), wanted)
}

func TestReportCaller(t *testing.T) {
	output := &bytes.Buffer{}

	log := Logger{
		MinLevel:     "verbose",
		ReportCaller: true,
		Output:       output,
	}

	log.Debug("Hello there.")

	wanted := `logger_test.go:77`
	assert.Contains(t, output.String(), wanted)
}

func TestReportCallerJSON(t *testing.T) {
	output := &bytes.Buffer{}

	log := Logger{
		MinLevel:     "verbose",
		ReportCaller: true,
		UseJSON:      true,
		Output:       output,
	}

	log.Debug("Hello there.")

	wanted := `logger_test.go:93`
	assert.Contains(t, output.String(), wanted)
}

func TestWithColor(t *testing.T) {
	output := &bytes.Buffer{}

	log := Logger{
		MinLevel: "verbose",
		Output:   output,
		UseColor: true,
	}

	log.Info("Hello there.")

	wanted := "\x1b[38;5;38mINFO\x1b[0m Hello there"
	assert.Contains(t, output.String(), wanted)
}

func TestWithTimeFormat(t *testing.T) {
	output := &bytes.Buffer{}

	log := Logger{
		MinLevel:   "verbose",
		Output:     output,
		TimeFormat: "15:04",
	}

	log.Info("Hello there.")

	wanted := time.Now().Format("15:04")
	assert.Contains(t, output.String(), wanted)
}

func TestLogger_Info(t *testing.T) {
	output := &bytes.Buffer{}

	log := Logger{
		MinLevel: "verbose",
		Output:   output,
	}

	log.Info("Hello there.")

	wanted := `INFO Hello there`
	assert.Contains(t, output.String(), wanted)
}

func TestLogger_Warning(t *testing.T) {
	output := &bytes.Buffer{}

	log := Logger{
		MinLevel: "verbose",
		Output:   output,
	}

	log.Warning("Hello there.")

	wanted := `WARNING Hello there`
	assert.Contains(t, output.String(), wanted)
}

func TestLogger_Verbose(t *testing.T) {
	output := &bytes.Buffer{}

	log := Logger{
		MinLevel: "verbose",
		Output:   output,
	}

	log.Verbose("Hello there.")

	wanted := `VERBOSE Hello there`
	assert.Contains(t, output.String(), wanted)
}
