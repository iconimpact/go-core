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

// Package logger Colorful logging with additional time and level output
package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

// Logger settings
type Logger struct {
	MinLevel     string
	TimeFormat   string // use example date: "2006-01-02 15:04:05.000" or "15:04:05.000"
	UseColor     bool
	UseJSON      bool
	ReportCaller bool
	Output       io.Writer
}

// LogRow represents a log entry and is serialized to JSON.
type LogRow struct {
	Level    string `json:"level"`
	Time     string `json:"time"`
	Message  string `json:"message"`
	LevelInt int    `json:"level_int"`
	Caller   string `json:"caller,omitempty"`
}

var logLevels = map[string]int{
	"verbose": 0,
	"debug":   1,
	"info":    2,
	"warning": 3,
	"error":   4,
}

// Colors ANSI CLI colors
type Colors struct {
	reset  string
	escape string
	silver string
	green  string
	blue   string
	yellow string
	red    string
}

func (l *Logger) formatJSON(level string, text ...interface{}) string {
	var timeText string

	logRow := LogRow{}
	if timeText = time.Now().Format(time.RFC3339); timeText != "" {
		logRow.Time = timeText
	}
	//logRow.Service = l.Service
	logRow.Level = level
	logRow.Caller = l.caller()

	formattedText := fmt.Sprintf("%v", text)
	formattedText = strings.Trim(formattedText, "[[")
	formattedText = strings.Trim(formattedText, "]]")

	logRow.Message = formattedText
	logRow.LevelInt = logLevels[level]

	encoded, err := json.Marshal(logRow)

	if err != nil {
		encoded, _ = json.Marshal(map[string]string{
			"time":    timeText,
			"message": "Error marshaling JSON",
			"level":   "error",
		})
	}

	return string(encoded[:])

}

// adds timestamp and color
func (l *Logger) formatText(level string, text ...interface{}) string {
	// get color for level
	colors := Colors{}
	if l.UseColor {
		colors = Colors{
			escape: "\x1b[38;5;",
			reset:  "\x1b[0m",
			silver: "251m",
			green:  "35m",
			blue:   "38m",
			yellow: "178m",
			red:    "197m",
		}
	}
	var coloredLevel = colors.escape + colors.silver + "VERBOSE" + colors.reset

	switch level {
	case "debug":
		coloredLevel = colors.escape + colors.green + "DEBUG" + colors.reset
	case "info":
		coloredLevel = colors.escape + colors.blue + "INFO" + colors.reset
	case "warning":
		coloredLevel = colors.escape + colors.yellow + "WARNING" + colors.reset
	case "error":
		coloredLevel = colors.escape + colors.red + "ERROR" + colors.reset
	}

	var ret = ""
	if timeText := time.Now().Format(l.TimeFormat); timeText != "" {
		ret += fmt.Sprintf("%v ", timeText)
	}
	if l.ReportCaller {
		ret += fmt.Sprintf("caller: %s ", l.caller())
	}
	formattedText := fmt.Sprintf("%v", text)
	formattedText = strings.Trim(formattedText, "[[")
	formattedText = strings.Trim(formattedText, "]]")
	ret += fmt.Sprintf("%v %v", coloredLevel, formattedText)

	return ret
}

func (l *Logger) caller() string {
	if !l.ReportCaller {
		return ""
	}

	// stack frames to skip
	loggerFrames := 4
	_, file, line, ok := runtime.Caller(loggerFrames)
	if !ok {
		return ""
	}

	return fmt.Sprintf("%s:%d", getFileName(file), line)
}

func (l *Logger) log(level string, text ...interface{}) {
	if l.Output == nil {
		l.Output = os.Stderr
	}

	if l.UseJSON {
		fmt.Fprintln(l.Output, l.formatJSON(level, text))
		return
	}

	fmt.Fprintln(l.Output, l.formatText(level, text))
}

// Verbose logging - note: Its just debug.
func (l *Logger) Verbose(text ...interface{}) {
	if l.MinLevel == "" || l.MinLevel == "verbose" {
		l.log("verbose", text)
	}
}

// Debug logging
func (l *Logger) Debug(text ...interface{}) {
	if l.MinLevel == "" || l.MinLevel == "verbose" || l.MinLevel == "debug" {
		l.log("debug", text)
	}
}

// Info logging
func (l *Logger) Info(text ...interface{}) {
	if l.MinLevel == "" || l.MinLevel == "verbose" || l.MinLevel == "debug" || l.MinLevel == "info" {
		l.log("info", text)
	}
}

// Warning logging
func (l *Logger) Warning(text ...interface{}) {
	if l.MinLevel == "" || l.MinLevel == "verbose" || l.MinLevel == "debug" || l.MinLevel == "info" ||
		l.MinLevel == "warning" {
		l.log("warning", text)
	}
}

// Error logging
func (l *Logger) Error(text ...interface{}) {
	if l.MinLevel == "" || l.MinLevel == "verbose" || l.MinLevel == "debug" || l.MinLevel == "info" ||
		l.MinLevel == "warning" || l.MinLevel == "error" {
		l.log("error", text)
	}
}

func getFileName(file string) string {
	_, fileName := path.Split(file)
	return fileName
}
