package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sakuradon99/gokit/trace"
	"github.com/sakuradon99/ioc"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"path"
	"strings"
	"sync"
)

var defaultLogger *logrus.Logger

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

var loadOnce sync.Once

func load() {
	if defaultLogger != nil {
		return
	}

	loadOnce.Do(func() {
		logger, err := ioc.GetObject[logrus.Logger]("")
		if err != nil {
			panic(err)
		}

		defaultLogger = logger.(*logrus.Logger)
	})
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
	if config.RotationMaxSize == 0 {
		config.RotationMaxSize = 50
	}
	if config.RotationMaxBackups == 0 {
		config.RotationMaxBackups = 3
	}
	if config.RotationMaxAge == 0 {
		config.RotationMaxAge = 28
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

	fileDebug := &lumberjack.Logger{
		Filename:   path.Join(config.Path, "debug.log"),
		MaxSize:    config.RotationMaxSize,
		MaxBackups: config.RotationMaxBackups,
		MaxAge:     config.RotationMaxAge,
	}

	fileInfo := &lumberjack.Logger{
		Filename:   path.Join(config.Path, "info.log"),
		MaxSize:    config.RotationMaxSize,
		MaxBackups: config.RotationMaxBackups,
		MaxAge:     config.RotationMaxAge,
	}

	fileWarn := &lumberjack.Logger{
		Filename:   path.Join(config.Path, "warn.log"),
		MaxSize:    config.RotationMaxSize,
		MaxBackups: config.RotationMaxBackups,
		MaxAge:     config.RotationMaxAge,
	}

	fileError := &lumberjack.Logger{
		Filename:   path.Join(config.Path, "error.log"),
		MaxSize:    config.RotationMaxSize,
		MaxBackups: config.RotationMaxBackups,
		MaxAge:     config.RotationMaxAge,
	}

	fileData := &lumberjack.Logger{
		Filename:   path.Join(config.Path, "data.log"),
		MaxSize:    config.RotationMaxSize,
		MaxBackups: config.RotationMaxBackups,
		MaxAge:     config.RotationMaxAge,
	}

	fileAccess := &lumberjack.Logger{
		Filename:   path.Join(config.Path, "access.log"),
		MaxSize:    config.RotationMaxSize,
		MaxBackups: config.RotationMaxBackups,
		MaxAge:     config.RotationMaxAge,
	}

	logger.AddHook(&LevelFileHook{
		levelToWriter: map[logrus.Level]*lumberjack.Logger{
			debugLevel:  fileDebug,
			infoLevel:   fileInfo,
			warnLevel:   fileWarn,
			errorLevel:  fileError,
			dataLevel:   fileData,
			accessLevel: fileAccess,
		},
		logFormatter: formatter,
	})

	return logger
}

type LevelFileHook struct {
	levelToWriter map[logrus.Level]*lumberjack.Logger
	logFormatter  logrus.Formatter
}

func (hook *LevelFileHook) Fire(entry *logrus.Entry) error {
	writer, ok := hook.levelToWriter[entry.Level]
	if !ok {
		return fmt.Errorf("no log writer for level: %v", entry.Level)
	}

	formattedLog, err := hook.logFormatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = writer.Write(formattedLog)
	return err
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
	load()
	defaultLogger.WithContext(ctx).WithFields(getLogFields(kv)).Log(debugLevel, message)
}

func Info(ctx context.Context, message string, kv ...LogField) {
	load()
	defaultLogger.WithContext(ctx).WithFields(getLogFields(kv)).Log(infoLevel, message)
}

func Warn(ctx context.Context, message string, kv ...LogField) {
	load()
	defaultLogger.WithContext(ctx).WithFields(getLogFields(kv)).Log(warnLevel, message)
}

func Error(ctx context.Context, message string, kv ...LogField) {
	load()
	defaultLogger.WithContext(ctx).WithFields(getLogFields(kv)).Log(errorLevel, message)
}

func Data(ctx context.Context, kv ...LogField) {
	load()
	defaultLogger.WithContext(ctx).WithFields(getLogFields(kv)).Log(dataLevel)
}

func Access(ctx context.Context, kv ...LogField) {
	load()
	defaultLogger.WithContext(ctx).WithFields(getLogFields(kv)).Log(accessLevel)
}
