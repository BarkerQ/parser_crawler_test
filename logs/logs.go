package logs

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"os"
)

// NewLogger Метод логирования действий
func NewLogger() log.Logger {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	logger = level.NewFilter(logger, level.AllowDebug())
	logger = log.With(logger,
		"date_error", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)

	return logger
}
