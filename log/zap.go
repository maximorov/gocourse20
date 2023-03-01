package log

import (
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

const customWriterKeyPrefix = `customwriter`

var (
	registerGcpOnce            = new(sync.Once)
	_               Optionable = &zapConfiger{}
	_               Logger     = &ZapExtended{}
)

type (
	//  Open-Closed Principle
	zapConfiger struct {
		zc *zap.Config
	}

	ZapExtended struct {
		*otelzap.SugaredLogger
	}
)

func (zo *zapConfiger) SetEnvironment(env EnvEnum) {
	var zc zap.Config

	switch env {
	default:
		fallthrough
	case DevEnv:
		zc = zap.NewDevelopmentConfig()
		zc.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	case ProdEnv:
		zc = zap.NewProductionConfig()
	}

	zo.zc = &zc
}

func (zo *zapConfiger) SetServiceVersion(version string) {
	if zo.zc.InitialFields == nil {
		zo.zc.InitialFields = make(map[string]any)
	}
	zo.zc.InitialFields[`service.version`] = version
}

func (zo *zapConfiger) SetFormatter(encoding FormatterTypeEnum) {
	zo.zc.Encoding = regularToZapEncoding(encoding)
}

func (zo *zapConfiger) SetLevel(level LevelEnum) {
	zo.zc.Level = zap.NewAtomicLevelAt(regularToZapLevel(level))
}

func (zo *zapConfiger) SetEnableCaller(enable bool) {
	zo.zc.DisableCaller = !enable
}

func (zo *zapConfiger) SetEnableStacktrace(enable bool) {
	zo.zc.DisableStacktrace = !enable
}

func (zo *zapConfiger) SetOutput(output io.Writer) {
	switch output {
	case os.Stdout:
		zo.zc.OutputPaths = []string{`stdout`}
	case os.Stderr:
		zo.zc.OutputPaths = []string{`stderr`}
	default:
		atomic.AddUint64(&customOutputSufixCounter, 1)
		customWriterKey := fmt.Sprintf(`%s%s`, customWriterKeyPrefix, strconv.Itoa(int(customOutputSufixCounter)))
		_ = zap.RegisterSink(customWriterKey, func(u *url.URL) (zap.Sink, error) {
			return customOutput{output}, nil
		})

		customPath := fmt.Sprintf("%s:whatever", customWriterKey)
		zo.zc.OutputPaths = []string{customPath}
	}
}

func (zo *zapConfiger) build(opts ...Option) zap.Config {
	for i, o := range opts {
		funcName := runtime.FuncForPC(reflect.ValueOf(o).Pointer()).Name()
		if strings.Contains(funcName, `logger.WithEnvironment`) {
			o(zo)
			opts = append(opts[:i], opts[i+1:]...)
			goto opts
		}
	}
	zo.SetEnvironment(DevEnv)

opts:
	for _, o := range opts {
		o(zo)
	}

	return *zo.zc
}

// With adds a variadic number of fields to the logging context. It accepts a
// loosely-typed key-value pairs. The first element of the pair is used as the field key
// and the second as the field value.
//
// For example,
//
//	logger.With(
//	   "hello", "world",
//	   "failure", errors.New("oh no"),
//	   "count", 42,
//	   "user", User{Name: "alice"},
//	)
//
// Note that the keys in key-value pairs should be strings.
func (ze ZapExtended) With(args ...any) Logger {
	return &ZapExtended{ze.SugaredLogger.With(args...)}
}

// WithFields deprecated. Use With instead
func (ze ZapExtended) WithFields(atrs ...Attribute) Logger {
	var l Logger
	l = &ze

	for _, a := range atrs {
		l = l.With(a.Key, a.Value)
	}

	return l
}

func NewZap(opts ...Option) (*ZapExtended, error) {
	var err error
	registerGcpOnce.Do(func() {
		err = zap.RegisterEncoder("gcp", NewGcpEncoder)
	})
	if err != nil {
		return nil, err
	}

	cb := &zapConfiger{}
	zl, err := cb.build(opts...).Build()
	if err != nil {
		return nil, err
	}

	return &ZapExtended{otelzap.New(zl).Sugar()}, nil
}
