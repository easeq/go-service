package zap

import (
	goconfig "github.com/easeq/go-config"
	uber_zap "go.uber.org/zap"
)

type zap struct {
	initializer   interface{}
	config        *Config
	sugaredLogger *uber_zap.SugaredLogger
}

func NewZap() *zap {
	config := goconfig.NewEnvConfig(new(Config)).(*Config)
	// log.Println("config", config, config.ZapConfig())
	// logger, err := config.ZapConfig().Build()
	// log.Println("logger", logger, err)
	// if err != nil {
	// 	panic(err)
	// }

	// var encoder zapcore.Encoder
	// if config.Encoding == "console" {
	// 	encoder = zapcore.NewConsoleEncoder(config.ZapConfig().EncoderConfig)
	// } else {
	// 	encoder = zapcore.NewJSONEncoder(config.ZapConfig().EncoderConfig)
	// }

	// logWriter := zapcore.AddSync(os.Stdout)
	// core := zapcore.NewCore(encoder, logWriter, config.AtomicLevel())
	// logger := uber_zap.New(core, uber_zap.AddCaller(), uber_zap.AddCallerSkip(1))

	logger, _ := uber_zap.NewDevelopment()
	defer logger.Sync()

	sugaredLogger := logger.Sugar()
	// if err := sugaredLogger.Sync(); err != nil {
	// 	panic(err)
	// }

	return &zap{nil, config, sugaredLogger}
}

func (l *zap) Debug(args ...interface{}) {
	l.sugaredLogger.Debug(args...)
}

func (l *zap) Debugf(template string, args ...interface{}) {
	l.sugaredLogger.Debugf(template, args...)
}

func (l *zap) Debugw(message string, args ...interface{}) {
	l.sugaredLogger.Debugw(message, args...)
}

func (l *zap) Info(args ...interface{}) {
	l.sugaredLogger.Info(args...)
}

func (l *zap) Infof(template string, args ...interface{}) {
	l.sugaredLogger.Infof(template, args...)
}

func (l *zap) Infow(message string, args ...interface{}) {
	l.sugaredLogger.Infow(message, args...)
}

func (l *zap) Warn(args ...interface{}) {
	l.sugaredLogger.Warn(args...)
}

func (l *zap) Warnf(template string, args ...interface{}) {
	l.sugaredLogger.Warnf(template, args...)
}

func (l *zap) Warnw(message string, args ...interface{}) {
	l.sugaredLogger.Warnw(message, args...)
}

func (l *zap) Error(args ...interface{}) {
	l.sugaredLogger.Error(args...)
}

func (l *zap) Errorf(template string, args ...interface{}) {
	l.sugaredLogger.Errorf(template, args...)
}

func (l *zap) Errorw(message string, args ...interface{}) {
	l.sugaredLogger.Errorw(message, args...)
}

func (l *zap) DPanic(args ...interface{}) {
	l.sugaredLogger.DPanic(args...)
}

func (l *zap) DPanicf(template string, args ...interface{}) {
	l.sugaredLogger.DPanicf(template, args...)
}

func (l *zap) DPanicw(message string, args ...interface{}) {
	l.sugaredLogger.DPanicw(message, args...)
}

func (l *zap) Panic(args ...interface{}) {
	l.sugaredLogger.Panic(args...)
}

func (l *zap) Panicf(template string, args ...interface{}) {
	l.sugaredLogger.Panicf(template, args...)
}

func (l *zap) Panicw(message string, args ...interface{}) {
	l.sugaredLogger.Panicw(message, args...)
}

func (l *zap) Fatal(args ...interface{}) {
	l.sugaredLogger.Fatal(args...)
}

func (l *zap) Fatalf(template string, args ...interface{}) {
	l.sugaredLogger.Fatalf(template, args...)
}

func (l *zap) Fatalw(message string, args ...interface{}) {
	l.sugaredLogger.Fatalw(message, args...)
}

// // AddDependency adds necessary service components as dependencies
// func (l *zap) AddDependency(dep interface{}) error {
// 	return nil
// }

// // Dependencies returns the string names of service components
// // that are required as dependencies for this component
// func (l *zap) Dependencies() []string {
// 	return []string{}
// }

// // CanRun returns true if the component has anything to Run
// func (l *zap) CanRun() bool {
// 	return false
// }

// // Run start the service component
// func (l *zap) Run(ctx context.Context) error {
// 	return nil
// }
