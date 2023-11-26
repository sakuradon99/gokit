package logger

type Config struct {
	Path       string `value:"logger.path;optional"`
	MaxSize    int    `value:"logger.max_size;optional"`
	MaxAge     int    `value:"logger.max_age;optional"`
	MaxBackups int    `value:"logger.max_backups;optional"`
}
