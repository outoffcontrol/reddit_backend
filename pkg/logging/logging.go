package logging

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	"go.uber.org/zap"
)

// const (
// 	requestIDKey = "requestID"
// 	loggerKey    = "logger"
// )

type Logger struct {
	Zap          *zap.SugaredLogger
	RequestIDKey string
	LoggerKey    string
}

func (ac *Logger) SetupReqID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			// https://github.com/opentracing/specification/blob/master/rfc/trace_identifiers.md
			requestID = RandBytesHex(16)
			r.Header.Set("X-Request-ID", requestID)
			r.Header.Set("trace-id", requestID)
			w.Header().Set("trace-id", requestID)
			w.Header().Set("X-Request-ID", requestID)
		}
		ctx := context.WithValue(r.Context(), ac.RequestIDKey, requestID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (ac *Logger) SetupLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		minLevel := zap.DebugLevel
		ctxlogger := ac.Zap.With(
			zap.String("logger", "ctxlog"),
			zap.String("trace-id", ac.RequestIDFromContext(r.Context())),
			zap.String("email", "outoffcontrol77@gmail.com"),
		).WithOptions(
			zap.IncreaseLevel(minLevel),
			zap.AddCaller(),
			// zap.AddCallerSkip(1),
			zap.AddStacktrace(zap.ErrorLevel),
		)

		ctx := context.WithValue(r.Context(), ac.LoggerKey, ctxlogger)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// func (ac *Logger) SetupAccessLog(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		next.ServeHTTP(w, r)

// 		ac.Z(r.Context()).Infow(r.URL.Path,
// 			"method", r.Method,
// 			"remote_addr", r.RemoteAddr,
// 			"url", r.URL.Path,
// 			"work_time", time.Since(start),
// 		)
// 	})
// }

func (ac *Logger) Z(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return ac.Zap
	}
	zap, ok := ctx.Value(ac.LoggerKey).(*zap.SugaredLogger)
	if !ok || zap == nil {
		return ac.Zap
	}
	return zap
}

func (ac *Logger) RequestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(ac.RequestIDKey).(string)
	if !ok {
		return "-"
	}
	return requestID
}

func RandBytesHex(n int) string {
	return fmt.Sprintf("%x", RandBytes(n))
}

func RandBytes(n int) []byte {
	res := make([]byte, n)
	rand.Read(res)
	return res
}
