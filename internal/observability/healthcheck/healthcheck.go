package healthcheck

import (
	"context"

	"github.com/babylonchain/staking-api-service/internal/queue"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logger zerolog.Logger = log.Logger

func SetLogger(customLogger zerolog.Logger) {
	logger = customLogger
}

func StartHealthCheckCron(ctx context.Context, queues *queue.Queues, cronTime string) error {
	c := cron.New()
	logger.Info().Msg("Initiated Health Check Cron")

	if cronTime == "" {
		cronTime = "@every 1m"
	}

	_, err := c.AddFunc(cronTime, func() {
		queueHealthCheck(queues)
	})

	if err != nil {
		return err
	}

	c.Start()

	go func() {
		<-ctx.Done()
		logger.Info().Msg("Stopping Health Check Cron")
		c.Stop()
	}()

	return nil
}

func queueHealthCheck(queues *queue.Queues) {
	if err := queues.IsConnectionHealthy(); err != nil {
		logger.Error().Err(err).Msg("One or more queue connections are not healthy.")
	}
}
