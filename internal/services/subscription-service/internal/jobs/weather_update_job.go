package jobs

import (
	"context"
	"time"

	"subscription-service/internal/domain"
	"subscription-service/internal/repository/subscriptions"
	"subscription-service/internal/proto"
)

type subscriptionRepositoryManager interface {
	GetDueSubscriptions(ctx context.Context) ([]subscriptions.Subscription, error)
	UpdateNextNotification(ctx context.Context, subscriptionID int64, t time.Time) error
}

type eventPublisherManager interface {
	PublishWithTopic(ctx context.Context, topic string, event any) error
}

type weatherClientManager interface {
	GetWeather(ctx context.Context, req *proto.WeatherRequest) (*proto.WeatherResponse, error)
}

type loggerManager interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Sync() error
}

type WeatherUpdateJob struct {
	repo          subscriptionRepositoryManager
	publisher     eventPublisherManager
	weatherClient weatherClientManager
	logger loggerManager
}

func NewWeatherUpdateJob(
	repo subscriptionRepositoryManager,
	publisher eventPublisherManager,
	weatherClient weatherClientManager,
	logger loggerManager,
) *WeatherUpdateJob {
	return &WeatherUpdateJob{
		repo:          repo,
		publisher:     publisher,
		weatherClient: weatherClient,
		logger: logger,
	}
}

func (j *WeatherUpdateJob) Run(ctx context.Context) {
	subscriptions, err := j.repo.GetDueSubscriptions(ctx)
	if err != nil {
		j.logger.Error("[WeatherUpdateJob] failed to get due subscriptions: %v", err)
		return
	}
	for _, s := range subscriptions {
		weatherResp, err := j.weatherClient.GetWeather(ctx, &proto.WeatherRequest{City: s.City})
		if err != nil {
			j.logger.Error("[WeatherUpdateJob] failed to get weather for city=%s: %v", s.City, err)
			continue
		}
		event := domain.WeatherUpdateEvent{
			Email:       s.ChannelValue,
			Metrics: domain.WeatherMetrics{
				City:        s.City,
				Description: weatherResp.Description,
				Temperature: weatherResp.Temperature,
				Humidity:    weatherResp.Humidity,
			},
			UpdatedAt: time.Now().Unix(),
		}

		if err := j.publisher.PublishWithTopic(ctx, "weather.updated", event); err != nil {
			j.logger.Error("[WeatherUpdateJob] failed to publish weather update for user=%d: %v", s.ID, err)
		}

		if err := j.repo.UpdateNextNotification(ctx, int64(s.ID), time.Now()); err != nil {
			j.logger.Error("[WeatherUpdateJob] failed to update next notification for user=%d: %v", s.ID, err)
		}
	}
}

func (j *WeatherUpdateJob) StartPeriodic(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			j.Run(ctx)
		case <-ctx.Done():
			j.logger.Info("[WeatherUpdateJob] stopped")
			return
		}
	}
}
