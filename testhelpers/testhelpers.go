package testhelpers

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// ObserveLogs returns a new observable test logger and observed logs based on
// the specified test logger.
func ObserveLogs(
	t *testing.T,
	testLog *zap.SugaredLogger,
) (*zap.SugaredLogger, *observer.ObservedLogs) {

	inMemoryLogCore, observedLogs := observer.New(zapcore.DebugLevel)
	observableLog := zap.New(zapcore.NewTee(testLog.Desugar().Core(), inMemoryLogCore))
	return observableLog.Sugar(), observedLogs
}

// RequireLastLogEntry tests that last log entry from the observed logs
// has the expected level, it's message contains the expected message and/or
// has fields having values that contain the expected fields values.
func RequireLastLogEntry(
	t *testing.T,
	observedLogs *observer.ObservedLogs,
	expectedLevel zapcore.Level,
	expectedMessage string,
	expectedFieldsContaining map[string]string,
) {

	logs := observedLogs.All()

	require.NotEmpty(t, logs)
	lastLog := logs[len(logs)-1]
	require.Equal(t, expectedLevel, lastLog.Level)

	if len(expectedMessage) > 0 {
		require.Contains(t, lastLog.Message, expectedMessage)
	}

	if len(expectedFieldsContaining) > 0 {
		actualFields := make(map[string]string)
		for _, actualField := range lastLog.Context {
			actualFields[actualField.Key] = fmt.Sprintf("%s", actualField.Interface)
		}

		foundFields := make(map[string]string)
		for expectedField, expectedValueContains := range expectedFieldsContaining {
			foundValue, ok := actualFields[expectedField]
			if !ok {
				continue
			}
			if strings.Contains(foundValue, expectedValueContains) {
				foundFields[expectedField] = foundValue
			}
		}

		require.Len(
			t,
			foundFields,
			len(expectedFieldsContaining),
			"Expected fields containing: %#v.\nFound fields containing: %#v.\nAll actual fields: %#v",
			expectedFieldsContaining, foundFields, actualFields)
	}
}
