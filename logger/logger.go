package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sakuradon99/gokit/trace"
	"github.com/sakuradon99/ioc"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"path"
	"strings"
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

func load() {
	if defaultLogger != nil {
		return
	}

	logger, err := ioc.GetObject[logrus.Logger]("")
	if err != nil {
		panic(err)
	}

	defaultLogger = logger.(*logrus.Logger)
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
	sb.WriteByte('\n')

	return []byte(sb.String()), nil
}

func newLogger(config *Config) *logrus.Logger {
	logger := logrus.New()
	logFormatter := &logFormatter{}

	logger.SetLevel(debugLevel)
	logger.SetFormatter(logFormatter)

	if config.Path == "" {
		config.Path = "log"
	}
	if config.MaxSize == 0 {
		config.MaxSize = 50
	}
	if config.MaxBackups == 0 {
		config.MaxBackups = 3
	}
	if config.MaxAge == 0 {
		config.MaxAge = 28
	}

	fileDebug := &lumberjack.Logger{
		Filename:   path.Join(config.Path, "debug.log"),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
	}

	fileInfo := &lumberjack.Logger{
		Filename:   path.Join(config.Path, "info.log"),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
	}

	fileWarn := &lumberjack.Logger{
		Filename:   path.Join(config.Path, "warn.log"),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
	}

	fileError := &lumberjack.Logger{
		Filename:   path.Join(config.Path, "error.log"),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
	}

	fileData := &lumberjack.Logger{
		Filename:   path.Join(config.Path, "data.log"),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
	}

	fileAccess := &lumberjack.Logger{
		Filename:   path.Join(config.Path, "access.log"),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
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
		logFormatter: logFormatter,
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
