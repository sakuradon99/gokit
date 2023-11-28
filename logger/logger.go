package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rifflock/lfshook"
	"github.com/sakuradon99/gokit/trace"
	"github.com/sirupsen/logrus"
	"io"
	"path"
	"strings"
)

var defaultLogger = logrus.New()

const (
	infoLevel logrus.Level = iota + 100
	errorLevel
	warnLevel
	dataLevel
	accessLevel
	debugLevel
)

var levelNames = map[logrus.Level]string{
	infoLevel:   "INFO",
	errorLevel:  "ERROR",
	warnLevel:   "WARN",
	dataLevel:   "DATA",
	accessLevel: "ACCESS",
	debugLevel:  "DEBUG",
}

type LogField struct {
	Key   string
	Value any
}

func Field(key string, val any) LogField {
	if err, ok := val.(error); ok {
		return LogField{Key: key, Value: err.Error()}
	}
	return LogField{Key: key, Value: val}
}

type logFormatter struct {
	maxBytes int
}

func (l *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	sb := &strings.Builder{}
	sb.WriteString(entry.Time.Format("2006-01-02 15:04:05.000") + "|")
	sb.WriteString(levelNames[entry.Level] + "|")
	traceID := trace.GetTraceID(entry.Context)
	if traceID == "" {
		traceID = "-"
	}
	sb.WriteString(traceID)
	if entry.Level < dataLevel {
		sb.WriteString("|")
		sb.WriteString(entry.Message)
	}
	if len(entry.Data) > 0 {
		jsonBytes, _ := json.Marshal(entry.Data)
		sb.WriteString("|" + string(jsonBytes))
	}

	if sb.Len() > l.maxBytes {
		s := sb.String()
		sb = &strings.Builder{}
		sb.WriteString(s[:l.maxBytes])
		sb.WriteString(fmt.Sprintf("...(%d bytes truncated)", len(s)-l.maxBytes))
	}

	sb.WriteByte('\n')

	return []byte(sb.String()), nil
}

func newLogger(config *Config) *logrus.Logger {
	if config.Path == "" {
		config.Path = "log"
	}
	if config.MaxBytes == 0 {
		config.MaxBytes = 1024
	}

	logger := logrus.New()
	formatter := &logFormatter{
		maxBytes: config.MaxBytes,
	}

	logger.SetLevel(debugLevel)
	logger.SetFormatter(formatter)
	if !config.EnableTerminalLog {
		logger.SetOutput(io.Discard)
	}

	pathMap := lfshook.PathMap{
		debugLevel:  path.Join(config.Path, "debug.log"),
		infoLevel:   path.Join(config.Path, "info.log"),
		warnLevel:   path.Join(config.Path, "warn.log"),
		errorLevel:  path.Join(config.Path, "error.log"),
		dataLevel:   path.Join(config.Path, "data.log"),
		accessLevel: path.Join(config.Path, "access.log"),
	}

	logger.AddHook(&LevelFileHook{lfshook.NewHook(pathMap, formatter)})
	defaultLogger = logger

	return logger
}

type LevelFileHook struct {
	*lfshook.LfsHook
}

func (hook *LevelFileHook) Levels() []logrus.Level {
	return []logrus.Level{
		debugLevel,
		infoLevel,
		warnLevel,
		errorLevel,
		dataLevel,
		accessLevel,
	}
}

func getLogFields(kvs []LogField) logrus.Fields {
	fields := logrus.Fields{}
	for _, kv := range kvs {
		fields[kv.Key] = kv.Value
	}
	return fields
}

func Debug(ctx context.Context, message string, kv ...LogField) {
	defaultLogger.WithContext(ctx).WithFields(getLogFields(kv)).Log(debugLevel, message)
}

func Info(ctx context.Context, message string, kv ...LogField) {
	defaultLogger.WithContext(ctx).WithFields(getLogFields(kv)).Log(infoLevel, message)
}

func Warn(ctx context.Context, message string, kv ...LogField) {
	defaultLogger.WithContext(ctx).WithFields(getLogFields(kv)).Log(warnLevel, message)
}

func Error(ctx context.Context, message string, kv ...LogField) {
	defaultLogger.WithContext(ctx).WithFields(getLogFields(kv)).Log(errorLevel, message)
}

func Data(ctx context.Context, kv ...LogField) {
	defaultLogger.WithContext(ctx).WithFields(getLogFields(kv)).Log(dataLevel)
}

func Access(ctx context.Context, kv ...LogField) {
	defaultLogger.WithContext(ctx).WithFields(getLogFields(kv)).Log(accessLevel)
}
