package processor

import (
    "testing"
)

func TestDefaultLoggerLogsWarnings(t *testing.T) {
	logger := DefaultLogger()
	if logger.config.LevelToLog != WARNING {
		t.FailNow()
	}
}

