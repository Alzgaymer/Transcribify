package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type configModifierFunc func(*zap.Config)

func New(configModifiers ...configModifierFunc) (*zap.Logger, error) {

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
func WithLevel(level zap.AtomicLevel) configModifierFunc {
	return func(config *zap.Config) {
		config.Level = level
	}
}

// WithEncoding "console" as default
func WithEncoding(encoding string) configModifierFunc {
	return func(config *zap.Config) {
		config.Encoding = encoding
	}
}

// WithDevelopment false as default
func WithDevelopment(dev bool) configModifierFunc {
	return func(config *zap.Config) {
		config.Development = dev
	}
}

// WithOutputPaths "stderr" as default
func WithOutputPaths(paths []string) configModifierFunc {
	return func(config *zap.Config) {
		config.OutputPaths = paths
	}
}

// WithEncoderConfig as default NewDevelopment
func WithEncoderConfig(encoder zapcore.EncoderConfig) configModifierFunc {
	return func(config *zap.Config) {
		config.EncoderConfig = encoder
	}
}
