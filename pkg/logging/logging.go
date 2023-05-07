package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type optionsFunc func(*zap.Config)

func New(configModifiers ...optionsFunc) (*zap.Logger, error) {

	cfg := zap.Config{
		Encoding:         "console",
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		EncoderConfig:    zap.NewProductionEncoderConfig(),
	}

	for _, c := range configModifiers {
		c(&cfg)
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// WithLevel zap.InfoLevel as default
func WithLevel(level zap.AtomicLevel) optionsFunc {
	return func(config *zap.Config) {
		config.Level = level
	}
}

// WithEncoding "console" as default
func WithEncoding(encoding string) optionsFunc {
	return func(config *zap.Config) {
		config.Encoding = encoding
	}
}

// WithDevelopment false as default
func WithDevelopment(dev bool) optionsFunc {
	return func(config *zap.Config) {
		config.Development = dev
	}
}

// WithOutputPaths "stderr" as default
func WithOutputPaths(paths ...string) optionsFunc {
	return func(config *zap.Config) {
		config.OutputPaths = paths
	}
}

// WithEncoderConfig as default NewDevelopment
func WithEncoderConfig(encoder zapcore.EncoderConfig) optionsFunc {
	return func(config *zap.Config) {
		config.EncoderConfig = encoder
	}
}
