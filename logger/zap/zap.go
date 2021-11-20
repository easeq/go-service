package zap

import (
	"github.com/easeq/go-service/component"
	uber_zap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Zap struct {
	Config *Config
	Logger *uber_zap.SugaredLogger
}

func NewZap() *Zap {
	config := NewConfig()
	core := zapcore.NewCore(
		config.GetEncoder(),
		config.GetLogWriter(),
		config.AtomicLevel(),
	)
	logger := uber_zap.New(
		core,
		uber_zap.AddCaller(),
		uber_zap.AddCallerSkip(1),
		uber_zap.Fields(
			uber_zap.String("service", config.ServiceName),
		),
	)
	sugaredLogger := logger.Sugar()

	return &Zap{Config: config, Logger: sugaredLogger}
}

func (l *Zap) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}

func (l *Zap) Debugf(template string, args ...interface{}) {
	l.Logger.Debugf(template, args...)
}

func (l *Zap) Debugw(message string, args ...interface{}) {
	l.Logger.Debugw(message, args...)
}

func (l *Zap) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *Zap) Infof(template string, args ...interface{}) {
	l.Logger.Infof(template, args...)
}

func (l *Zap) Infow(message string, args ...interface{}) {
	l.Logger.Infow(message, args...)
}

func (l *Zap) Warn(args ...interface{}) {
	l.Logger.Warn(args...)
}

func (l *Zap) Warnf(template string, args ...interface{}) {
	l.Logger.Warnf(template, args...)
}

func (l *Zap) Warnw(message string, args ...interface{}) {
	l.Logger.Warnw(message, args...)
}

func (l *Zap) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

func (l *Zap) Errorf(template string, args ...interface{}) {
	l.Logger.Errorf(template, args...)
}

func (l *Zap) Errorw(message string, args ...interface{}) {
	l.Logger.Errorw(message, args...)
}

func (l *Zap) DPanic(args ...interface{}) {
	l.Logger.DPanic(args...)
}

func (l *Zap) DPanicf(template string, args ...interface{}) {
	l.Logger.DPanicf(template, args...)
}

func (l *Zap) DPanicw(message string, args ...interface{}) {
	l.Logger.DPanicw(message, args...)
}

func (l *Zap) Panic(args ...interface{}) {
	l.Logger.Panic(args...)
}

func (l *Zap) Panicf(template string, args ...interface{}) {
	l.Logger.Panicf(template, args...)
}

func (l *Zap) Panicw(message string, args ...interface{}) {
	l.Logger.Panicw(message, args...)
}

func (l *Zap) Fatal(args ...interface{}) {
	l.Logger.Fatal(args...)
}

func (l *Zap) Fatalf(template string, args ...interface{}) {
	l.Logger.Fatalf(template, args...)
}

func (l *Zap) Fatalw(message string, args ...interface{}) {
	l.Logger.Fatalw(message, args...)
}

func (l *Zap) HasInitializer() bool {
	return false
}

func (l *Zap) Initializer() component.Initializer {
	return nil
}
