package zap

import (
	"strings"

	goconfig "github.com/easeq/go-config"
	uber_zap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Netflix/go-env"
)

const (
	SEPARATOR = ","
)

// Config holds the etcd configuration
type Config struct {
	Level            string `env:"LOGGER_LEVEL,default=debug"`
	Encoding         string `env:"LOGGER_ENCODING,default=json"`
	OutputPaths      string `env:"LOGGER_OUTPUT_PATHS,default=stderr"`
	ErrorOutputPaths string `env:"LOGGER_ERROR_OUTPUT_PATHS,default=stderr"`
	EncConfig        EncoderConfig
}

type EncoderConfig struct {
	MessageKey       string `env:"LOGGER_ENC_MESSAGE_KEY,default=msg"`
	LevelKey         string `env:"LOGGER_ENC_LEVEL_KEY,default=level"`
	TimeKey          string `env:"LOGGER_ENC_TIME_KEY,default=ts"`
	NameKey          string `env:"LOGGER_ENC_NAME_KEY,default=logger"`
	CallerKey        string `env:"LOGGER_ENC_CALLER_KEY,default=caller"`
	FunctionKey      string `env:"LOGGER_ENC_FUNCTION_KEY,default="`
	StacktraceKey    string `env:"LOGGER_ENC_STACKTRACE_KEY,default=stacktrace"`
	LineEnding       string `env:"LOGGER_ENC_LINE_ENDING,default=\n"`
	ConsoleSeparator string `env:"LOGGER_ENC_CONSOLE_SEPARATOR,default=\t"`
}

// NewConfig returns the env config for etcd client
func NewConfig() *Config {
	return goconfig.NewEnvConfig(new(Config)).(*Config)
}

// UnmarshalEnv env.EnvSet to Config
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}

func (c *Config) AtomicLevel() uber_zap.AtomicLevel {
	al := uber_zap.AtomicLevel{}
	al.UnmarshalText([]byte(c.Level))
	return al
}

func (c *Config) OutputPathsArray() []string {
	if c.OutputPaths == "" {
		return []string{}
	}

	return strings.Split(c.OutputPaths, SEPARATOR)
}

func (c *Config) ErrorOutputPathsArray() []string {
	if c.ErrorOutputPaths == "" {
		return []string{}
	}

	return strings.Split(c.ErrorOutputPaths, SEPARATOR)
}

func (c *Config) EncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:       c.EncConfig.MessageKey,
		LevelKey:         c.EncConfig.LevelKey,
		TimeKey:          c.EncConfig.TimeKey,
		NameKey:          c.EncConfig.NameKey,
		CallerKey:        c.EncConfig.CallerKey,
		FunctionKey:      c.EncConfig.FunctionKey,
		StacktraceKey:    c.EncConfig.StacktraceKey,
		LineEnding:       c.EncConfig.LineEnding,
		ConsoleSeparator: c.EncConfig.ConsoleSeparator,
		EncodeTime:       zapcore.ISO8601TimeEncoder,
	}
}

func (c *Config) ZapConfig() *uber_zap.Config {
	cfg := &uber_zap.Config{
		Level:            c.AtomicLevel(),
		Encoding:         c.Encoding,
		OutputPaths:      c.OutputPathsArray(),
		ErrorOutputPaths: c.ErrorOutputPathsArray(),
		EncoderConfig:    c.EncoderConfig(),
	}

	return cfg
}
