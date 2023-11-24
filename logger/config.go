package logger

type Config struct {
	Path       string `value:"logger.path;optional"`
	MaxSize    int    `value:"logger.maxSize;optional"`
	MaxAge     int    `value:"logger.maxAge;optional"`
	MaxBackups int    `value:"logger.maxBackups;optional"`
}
