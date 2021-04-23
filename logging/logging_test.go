package logging

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdOutLogger(t *testing.T) {
	builder := strings.Builder{}
	verify := func(expected string) {
		actual := builder.String()
		assert.Equal(t, expected, actual)
		builder.Reset()
	}

	logger := NewStdOutLogger()
	stdLogger := logger.(*stdOutLogger)
	stdLogger.l.SetOutput(&builder)
	stdLogger.l.SetFlags(0) // disable time output

	logger.Info(nil, "without args")
	verify("info: without args \n")

	logger.Info(nil, "one arg", NewArg("arg1", "val1"))
	verify("info: one arg: arg1: val1\n")

	logger.Info(nil, "two args", NewArg("arg1", "val1"), NewArg("arg2", "val2"))
	verify("info: two args: arg1: val1, arg2: val2\n")

}
