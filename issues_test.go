package slog_test

import (
	"testing"
	"time"

	"github.com/tomorrowsky/slog"
	"github.com/tomorrowsky/slog/handler"
)

func TestIssues_27(t *testing.T) {
	defer slog.Reset()

	count := 0
	for {
		if count >= 6 {
			break
		}
		slog.Infof("info log %d", count)
		time.Sleep(time.Second)
		count++
	}
}

// https://github.com/gookit/slog/issues/31
func TestIssues_31(t *testing.T) {
	defer slog.Reset()
	defer slog.Flush()

	h1 := handler.MustFileHandler("testdata/error.log", true)
	// slog.DangerLevels equals slog.Levels{slog.PanicLevel, slog.PanicLevel, slog.ErrorLevel, slog.WarnLevel}
	h1.Levels = slog.DangerLevels

	h2 := handler.MustFileHandler("testdata/info.log", true)
	h2.Levels = slog.Levels{slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel, slog.TraceLevel}

	slog.PushHandler(h1)
	slog.PushHandler(h2)

	// add logs
	slog.Info("info message text")
	slog.Error("error message text")
}
