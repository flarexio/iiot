package iiot

import (
	"context"

	"go.uber.org/zap"
)

func LoggingMiddleware(log *zap.Logger) ServiceMiddleware {
	return func(next Service) Service {
		log := log.With(
			zap.String("service", "iiot"),
		)

		log.Info("service running")

		return &loggingMiddleware{
			log:  log,
			next: next,
		}
	}
}

type loggingMiddleware struct {
	log  *zap.Logger
	next Service
}

func (mw *loggingMiddleware) CheckConnection(ctx context.Context, network string, address string) error {
	log := mw.log.With(
		zap.String("action", "check_connection"),
		zap.String("network", network),
		zap.String("address", address),
	)

	err := mw.next.CheckConnection(ctx, network, address)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("Connection successful")
	return nil
}
