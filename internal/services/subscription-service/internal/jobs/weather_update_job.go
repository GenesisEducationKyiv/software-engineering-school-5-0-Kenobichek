package jobs

import (
	"context"
	"log"
	"time"

	"subscription-service/internal/domain"
	"subscription-service/internal/infrastructure"
	"subscription-service/internal/repository"
	"subscription-service/internal/proto"
)

type WeatherUpdateJob struct {
	repo          *repository.Repository
	publisher     *infrastructure.KafkaPublisher
	weatherClient proto.WeatherServiceClient
}

func NewWeatherUpdateJob(
	repo *repository.Repository,
	publisher *infrastructure.KafkaPublisher,
	weatherClient proto.WeatherServiceClient,
) *WeatherUpdateJob {
	return &WeatherUpdateJob{
		repo:          repo,
		publisher:     publisher,
		weatherClient: weatherClient,
	}
}

func (j *WeatherUpdateJob) Run(ctx context.Context) {
	subscriptions := j.repo.GetDueSubscriptions(ctx)

	for _, s := range subscriptions {
		weatherResp, err := j.weatherClient.GetWeather(ctx, &proto.WeatherRequest{City: s.City})
		if err != nil {
			log.Printf("[WeatherUpdateJob] failed to get weather for city=%s: %v", s.City, err)
			continue
		}
		
		// updated_at := time.Now().Unix()

		event := domain.WeatherUpdateEvent{
			Email:       s.ChannelValue,
			Metrics: domain.WeatherMetrics{
				City:        s.City,
				Description: weatherResp.Description,
				Temperature: weatherResp.Temperature,
				Humidity:    weatherResp.Humidity,
			},
			UpdatedAt:   time.Now().Unix(),
		}

		if err := j.publisher.PublishWithTopic(ctx, "weather.updated", event); err != nil {
			log.Printf("[WeatherUpdateJob] failed to publish weather update for user=%d: %v", s.ID, err)
		}

		if err := j.repo.UpdateNextNotification(ctx, s.ID, time.Now()); err != nil {
			log.Printf("[WeatherUpdateJob] failed to update next notification for user=%d: %v", s.ID, err)
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
			log.Println("WeatherUpdateJob stopped")
			return
		}
	}
}
