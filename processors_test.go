package slog_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tomorrowsky/slog"
)

func TestLogger_AddProcessor(t *testing.T) {
	buf := new(bytes.Buffer)

	l := slog.NewJSONSugared(buf, slog.InfoLevel)
	l.AddProcessor(slog.AddHostname())
	l.Info("message")

	hostname, _ := os.Hostname()

	// {"channel":"application","data":{},"datetime":"2020/07/17 12:01:35","extra":{},"hostname":"InhereMac","level":"INFO","message":"message"}
	str := buf.String()
	buf.Reset()
	assert.Contains(t, str, `"level":"INFO"`)
	assert.Contains(t, str, `"message":"message"`)
	assert.Contains(t, str, fmt.Sprintf(`"hostname":"%s"`, hostname))

	l.ResetProcessors()
	l.PushProcessor(slog.MemoryUsage)
	l.Info("message2")

	// {"channel":"application","data":{},"datetime":"2020/07/16 16:40:18","extra":{"memoryUsage":326072},"level":"INFO","message":"message2"}
	str = buf.String()
	buf.Reset()
	assert.Contains(t, str, `"message":"message2"`)
	assert.Contains(t, str, `"memoryUsage":`)

	l.ResetProcessors()
	l.PushProcessor(slog.AddUniqueID("requestId"))
	l.Info("message3")
	str = buf.String()
	buf.Reset()
	assert.Contains(t, str, `"message":"message3"`)
	assert.Contains(t, str, `"requestId":`)
	fmt.Println(str)
}
