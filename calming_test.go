package processor

import (
	"testing"
)

func TestCalming(t *testing.T) {
	defer Calm()
	panic("Should not be displayed.")
}

func TestCalmingAndLog(t *testing.T) {
	defer CalmAndLog()
	panic("Should be logged.")
}

func TestCalmingAndLogFunc(t *testing.T) {
	defer CalmAndLogFunc("TestCalmingAndLogFunc")()
	panic("Should be logged with function name.")
}
