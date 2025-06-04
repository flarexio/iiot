package iiot

import (
	"context"
	"encoding/json"

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

func (mw *loggingMiddleware) ListDrivers(ctx context.Context) ([]string, error) {
	log := mw.log.With(
		zap.String("action", "list_drivers"),
	)

	drivers, err := mw.next.ListDrivers(ctx)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Info("drivers retrieved", zap.Int("count", len(drivers)))
	return drivers, nil
}

func (mw *loggingMiddleware) Schema(ctx context.Context, driver string) (json.RawMessage, error) {
	log := mw.log.With(
		zap.String("action", "schema"),
		zap.String("driver", driver),
	)

	schema, err := mw.next.Schema(ctx, driver)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Info("schema retrieved")
	return schema, nil
}

func (mw *loggingMiddleware) Instruction(ctx context.Context, driver string) (string, error) {
	log := mw.log.With(
		zap.String("action", "instruction"),
		zap.String("driver", driver),
	)

	instruction, err := mw.next.Instruction(ctx, driver)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	log.Info("instruction retrieved")
	return instruction, nil
}

func (mw *loggingMiddleware) ReadPoints(ctx context.Context, driver string, raw json.RawMessage) ([]any, error) {
	log := mw.log.With(
		zap.String("action", "read_points"),
		zap.String("driver", driver),
	)

	points, err := mw.next.ReadPoints(ctx, driver, raw)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Info("Read points successful", zap.Any("points", points))
	return points, nil
}
