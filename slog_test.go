package slog_test

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/gookit/goutil/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/tomorrowsky/slog"
	"github.com/tomorrowsky/slog/handler"
)

var doNothing = func(code int) {
	// do nothing
}

func TestStd(t *testing.T) {
	defer slog.Reset()
	assert.Equal(t, "stdLogger", slog.Std().Name())

	_, ok := slog.GetFormatter().(*slog.TextFormatter)
	assert.True(t, ok)

	slog.SetLogLevel(slog.WarnLevel)
	slog.SetFormatter(slog.NewJSONFormatter())

	assert.True(t, slog.Std().IsHandling(slog.WarnLevel))
	assert.True(t, slog.Std().IsHandling(slog.ErrorLevel))
	assert.False(t, slog.Std().IsHandling(slog.InfoLevel))

	_, ok = slog.GetFormatter().(*slog.JSONFormatter)
	assert.True(t, ok)

	buf := new(bytes.Buffer)
	slog.Std().ExitFunc = func(code int) {
		buf.WriteString("Exited,")
		buf.WriteString(strconv.Itoa(code))
	}
	slog.Exit(34)
	assert.Equal(t, "Exited,34", buf.String())
}

func TestTextFormatNoColor(t *testing.T) {
	defer slog.Reset()
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = false

		logger.ExitFunc = slog.DoNothingOnExit
	})

	printLogs("print log message")
	printfLogs("print log with %s", "params")

	slog.Reset()
}

type logTest struct {
	*slog.SugaredLogger
}

func (l logTest) testPrint() {
	l.Logger.Info("print testing")
}

func TestTextFormatWithColor(t *testing.T) {
	defer slog.Reset()

	slog.Configure(func(l *slog.SugaredLogger) {
		l.Level = slog.PanicLevel
		l.ExitFunc = slog.DoNothingOnExit
	})

	printLogs("this is a simple log message")
	slog.Std().Trace("this is a simple log message")

	lt := &logTest{
		slog.Std(),
	}
	lt.testPrint()

	slog.GetFormatter().(*slog.TextFormatter).SetTemplate(slog.NamedTemplate)
	printfLogs("print log with %s", "params")

	tpl := "[{{datetime}}] [{{channel}}] [{{level}}] [{{func}}] {{message}} {{data}} {{extra}}\n"
	slog.GetFormatter().(*slog.TextFormatter).SetTemplate(tpl)
	printfLogs("print log with %s", "params")

	lt = &logTest{
		slog.Std(),
	}
	lt.testPrint()
}

func printLogs(msg string) {
	slog.Print(msg)
	slog.Println(msg)
	slog.Trace(msg)
	slog.Debug(msg)
	slog.Info(msg)
	slog.Notice(msg)
	slog.Warn(msg)
	slog.Error(msg)
	slog.Fatal(msg)
	slog.ErrorT(errors.New(msg))
}

func printfLogs(msg string, args ...interface{}) {
	slog.Printf(msg, args...)
	slog.Tracef(msg, args...)
	slog.Debugf(msg, args...)
	slog.Infof(msg, args...)
	slog.Noticef(msg, args...)
	slog.Warnf(msg, args...)
	slog.Errorf(msg, args...)
	slog.Fatalf(msg, args...)
}

func TestUseJSONFormat(t *testing.T) {
	defer slog.Reset()
	slog.SetFormatter(slog.NewJSONFormatter())

	slog.Info("info log message")
	slog.Warn("warning log message")
	slog.WithData(slog.M{
		"key0": 134,
		"key1": "abc",
	}).Infof("info log %s", "message")

	r := slog.WithFields(slog.M{
		"category": "service",
		"IP":       "127.0.0.1",
	})
	r.Infof("info %s", "message")
	r.Debugf("debug %s", "message")
}

func TestAddHandler(t *testing.T) {
	defer slog.Reset()
	slog.AddHandler(handler.NewConsoleHandler(slog.AllLevels))

	h2 := handler.NewConsoleHandler(slog.AllLevels)
	h2.SetFormatter(slog.NewJSONFormatter().Configure(func(f *slog.JSONFormatter) {
		f.Aliases = slog.StringMap{
			"level":   "levelName",
			"message": "msg",
			"data":    "params",
		}
	}))

	slog.AddHandlers(h2)
	slog.Infof("info %s", "message")
}

func TestAddProcessor(t *testing.T) {
	defer slog.Reset()

	buf := new(bytes.Buffer)
	slog.Configure(func(logger *slog.SugaredLogger) {
		logger.Level = slog.TraceLevel
		logger.Output = buf
		logger.Formatter = slog.NewJSONFormatter()
	})

	slog.AddProcessor(slog.AddHostname())
	slog.Trace("Trace message")
	slog.Tracef("Tracef %s", "message")

	str := buf.String()
	buf.Reset()
	fmt.Println(str)
	assert.Contains(t, str, `"hostname":`)
	assert.Contains(t, str, "Trace message")
	assert.Contains(t, str, "Tracef message")

	slog.AddProcessors(slog.ProcessorFunc(func(r *slog.Record) {
		r.AddField("newField", "newValue")
	}))
	slog.Debug("Debug message")
	slog.Debugf("Debugf %s", "message")
	str = buf.String()
	buf.Reset()

	assert.Contains(t, str, `"newField":"newValue"`)
	assert.Contains(t, str, "Debug message")
	assert.Contains(t, str, "Debugf message")
}

func TestLevelName(t *testing.T) {
	for level, wantName := range slog.LevelNames {
		realName := slog.LevelName(level)
		assert.Equal(t, wantName, realName)
	}

	assert.Equal(t, "UNKNOWN", slog.LevelName(20))
}

func TestName2Level(t *testing.T) {
	for wantLevel, name := range slog.LevelNames {
		level, err := slog.Name2Level(name)
		assert.NoError(t, err)
		assert.Equal(t, wantLevel, level)
	}

	// special names
	tests := map[slog.Level]string{
		slog.WarnLevel:  "warn",
		slog.ErrorLevel: "err",
		slog.InfoLevel:  "",
	}
	for wantLevel, name := range tests {
		level, err := slog.Name2Level(name)
		assert.NoError(t, err)
		assert.Equal(t, wantLevel, level)
	}

	level, err := slog.Name2Level("unknown")
	assert.Error(t, err)
	assert.Equal(t, slog.Level(0), level)
}

func TestPrependExitHandler(t *testing.T) {
	defer slog.Reset()

	assert.Len(t, slog.ExitHandlers(), 0)

	buf := new(bytes.Buffer)
	slog.PrependExitHandler(func() {
		buf.WriteString("HANDLER1-")
	})
	slog.PrependExitHandler(func() {
		buf.WriteString("HANDLER2-")
	})
	assert.Len(t, slog.ExitHandlers(), 2)

	slog.SetExitFunc(func(code int) {
		buf.WriteString("Exited")
	})
	slog.Exit(23)
	assert.Equal(t, "HANDLER2-HANDLER1-Exited", buf.String())
}

func TestRegisterExitHandler(t *testing.T) {
	defer slog.Reset()

	assert.Len(t, slog.ExitHandlers(), 0)

	buf := new(bytes.Buffer)
	slog.RegisterExitHandler(func() {
		buf.WriteString("HANDLER1-")
	})
	slog.RegisterExitHandler(func() {
		buf.WriteString("HANDLER2-")
	})
	// prepend
	slog.PrependExitHandler(func() {
		buf.WriteString("HANDLER3-")
	})
	assert.Len(t, slog.ExitHandlers(), 3)

	slog.SetExitFunc(func(code int) {
		buf.WriteString("Exited")
	})
	slog.Exit(23)
	assert.Equal(t, "HANDLER3-HANDLER1-HANDLER2-Exited", buf.String())
}

func TestExitHandlerWithError(t *testing.T) {
	defer slog.Reset()

	assert.Len(t, slog.ExitHandlers(), 0)

	slog.RegisterExitHandler(func() {
		panic("test error")
	})

	slog.SetExitFunc(func(code int) {})

	testutil.RewriteStderr()
	slog.Exit(23)
	str := testutil.RestoreStderr()
	assert.Equal(t, "Run exit handler error: test error\n", str)
}

func TestLogger_ExitHandlerWithError(t *testing.T) {
	l := slog.NewWithConfig(func(l *slog.Logger) {
		l.ExitFunc = doNothing
	})

	assert.Len(t, l.ExitHandlers(), 0)

	l.RegisterExitHandler(func() {
		panic("test error")
	})

	testutil.RewriteStderr()
	l.Exit(23)
	str := testutil.RestoreStderr()
	assert.Equal(t, "Run exit handler error: test error\n", str)
}
