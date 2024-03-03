package zinc

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Logger interface {
	LogQuery(ctx context.Context, msg string)
	LogQueryErr(ctx context.Context, msg string)
}

// StdLogger

type stdLogger struct{}

var StdLogger Logger = stdLogger{}

func (stdLogger) LogQuery(_ context.Context, msg string) {
	log.Println(msg)
}

func (stdLogger) LogQueryErr(_ context.Context, msg string) {
	log.Println(msg)
}

// LogFormatter

type LogFormatter func(q string, args UnitedArgs, err error, elapsed int64) string

func defaultLogFormatter(q string, args UnitedArgs, err error, elapsed int64) string {
	return fmt.Sprintf("%s %dms", q, elapsed)
}

// util

func logDo[R any](
	ctx context.Context,
	q string, uArgs UnitedArgs,
	bound string, boundArgs []any,
	opts *Options,
	f func() (R, error),
) (R, error) {
	l := opts.Logger
	if l == nil {
		return f()
	}

	logFormatter := opts.LogFormatter
	logBound := opts.LogBound
	logSuccess := opts.LogSuccess
	slowThreshold := opts.LogSlowThreshold
	if logFormatter == nil {
		logFormatter = defaultLogFormatter
	}
	if slowThreshold <= 0 {
		slowThreshold = 3000
	}
	formatLog := func(err error, elapsed int64) string {
		if logBound {
			return logFormatter(bound, UnitedArgs{Unnamed: boundArgs}, err, elapsed)
		} else {
			return logFormatter(q, uArgs, err, elapsed)
		}
	}
	startAt := time.Now()
	sqlRes, err := f()
	elapsed := time.Since(startAt).Milliseconds()
	if err != nil {
		l.LogQueryErr(ctx, formatLog(err, elapsed))
	} else {
		if logSuccess || elapsed >= slowThreshold {
			l.LogQuery(ctx, formatLog(err, elapsed))
		}
	}
	return sqlRes, err
}
