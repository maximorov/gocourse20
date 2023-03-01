package log

import "io"

type Option = func(Optionable)
type Optionable interface {
	SetServiceVersion(string)
	SetFormatter(FormatterTypeEnum)
	SetLevel(LevelEnum)
	SetEnableCaller(bool)
	SetOutput(io.Writer)
	SetEnvironment(EnvEnum)
	SetEnableStacktrace(bool)
	DoSometing()
}

var DefaultOptions = []Option{
	WithEnvironment(DevEnv),
	WithFormatterType(DefaultFormatterType),
	WithLevel(InfoLevel),
	WithEnableCaller(false),
}

func WithFormatterType(formatterType FormatterTypeEnum) Option {
	return func(o Optionable) {
		o.SetFormatter(formatterType)
	}
}

func WithLevel(level LevelEnum) Option {
	return func(o Optionable) {
		o.SetLevel(level)
	}
}

func WithEnableCaller(enable bool) Option {
	return func(o Optionable) {
		o.SetEnableCaller(enable)
	}
}

func WithEnableStacktrace(enable bool) Option {
	return func(o Optionable) {
		o.SetEnableStacktrace(enable)
	}
}

func WithServiceVersion(serviceVersion string) Option {
	return func(o Optionable) {
		o.SetServiceVersion(serviceVersion)
	}
}

func WithEnvironment(env EnvEnum) Option {
	return func(o Optionable) {
		o.SetEnvironment(env)
	}
}

func WithOutput(output io.Writer) Option {
	return func(o Optionable) {
		o.SetOutput(output)
	}
}

func WithSomethingElse() Option {
	return func(o Optionable) {
		o.DoSometing()
	}
}
