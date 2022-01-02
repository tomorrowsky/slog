package handler

import (
	"github.com/tomorrowsky/slog"
)

// SimpleFileHandler struct
// - no buffer, will direct write logs to file.
type SimpleFileHandler struct {
	fileWrapper
	lockWrapper
	// LevelWithFormatter support level and formatter
	LevelWithFormatter
}

// MustSimpleFile new instance
func MustSimpleFile(filepath string) *SimpleFileHandler {
	h, err := NewSimpleFileHandler(filepath)
	if err != nil {
		panic(err)
	}

	return h
}

// NewSimpleFileHandler new instance
func NewSimpleFile(filepath string) (*SimpleFileHandler, error) {
	return NewSimpleFileHandler(filepath)
}

// NewSimpleFileHandler instance
//
// Usage:
// 	h, err := NewSimpleFileHandler("", DefaultFileFlags)
// custom file flags
// 	h, err := NewSimpleFileHandler("", os.O_CREATE | os.O_WRONLY | os.O_APPEND)
// custom formatter
//	h.SetFormatter(slog.NewJSONFormatter())
//	slog.PushHandler(h)
//	slog.Info("log message")
func NewSimpleFileHandler(filepath string) (*SimpleFileHandler, error) {
	fh := fileWrapper{fpath: filepath}
	if err := fh.ReopenFile(); err != nil {
		return nil, err
	}

	h := &SimpleFileHandler{
		fileWrapper: fh,
	}

	// init default log level
	h.Level = slog.InfoLevel

	return h, nil
}

// Handle the log record
func (h *SimpleFileHandler) Handle(r *slog.Record) (err error) {
	var bts []byte
	bts, err = h.Formatter().Format(r)
	if err != nil {
		return
	}

	// if enable lock
	h.Lock()
	defer h.Unlock()

	// direct write logs
	_, err = h.file.Write(bts)
	return
}
