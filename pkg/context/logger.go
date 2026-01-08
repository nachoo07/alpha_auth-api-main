package authcontext

import (
	"context"

	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

func NewLogger() {
	logger = log.New()
	logger.SetFormatter(&log.JSONFormatter{})
}

func Logger(ctx context.Context) *log.Entry {
	if ctx != nil {
		var requestID string
		newLogger := logger
		if ctxRqId, ok := ctx.Value(requestIDKey).(string); ok {
			requestID = ctxRqId
		} else {
			requestID = newRequestID()
		}

		return newLogger.WithFields(log.Fields{
			string(requestIDKey): requestID,
		})
	}

	return defaultLogger()
}

func defaultLogger() *log.Entry {
	newLogger := logger

	return newLogger.WithFields(log.Fields{
		string(requestIDKey): newRequestID(),
	})
}
