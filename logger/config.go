package logger

type Config struct {
	Path              string `value:"logger.path;optional"`
	EnableTerminalLog bool   `value:"logger.enable_terminal_log;optional"`
	MaxBytes          int    `value:"logger.max_bytes;optional"`
}
