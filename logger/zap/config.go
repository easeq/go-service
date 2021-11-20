package zap

import (
	"github.com/easeq/go-service/component"
	"github.com/natefinch/lumberjack"
	uber_zap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Netflix/go-env"
)

// Config holds the etcd configuration
type Config struct {
	ServiceName      string `env:"SERVICE_NAME"`
	Dev              bool   `env:"LOGGER_DEV_MODE,default=true"`
	Level            string `env:"LOGGER_LEVEL,default=debug"`
	Encoding         string `env:"LOGGER_ENCODING,default=json"`
	OutputPath       string `env:"LOGGER_OUTPUT_PATH,default=./data/service.log"`
	MaxFileSize      int    `env:"LOGGER_MAX_FILE_SIZE,default=10"`
	MaxNumBackups    int    `env:"LOGGER_MAX_NUM_BACKUPS,default=5"`
	MaxRetentionDays int    `env:"LOGGER_MAX_RETENTION_DAYS,default=30"`
	CompressOld      bool   `env:"LOGGER_COMPRESS_OLD,default=true"`
}

// NewConfig returns the parsed config for zap from env
func NewConfig() *Config {
	c := new(Config)
	component.NewConfig(c)

	return c
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

func (c *Config) GetLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   c.OutputPath,
		MaxSize:    c.MaxFileSize,
		MaxBackups: c.MaxNumBackups,
		MaxAge:     c.MaxRetentionDays,
		Compress:   c.CompressOld,
	}

	return zapcore.AddSync(lumberJackLogger)
}

func (c *Config) GetEncoder() zapcore.Encoder {
	var encoderConfig zapcore.EncoderConfig
	if c.Dev {
		encoderConfig = uber_zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = uber_zap.NewProductionEncoderConfig()
	}

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	if c.Encoding == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}
