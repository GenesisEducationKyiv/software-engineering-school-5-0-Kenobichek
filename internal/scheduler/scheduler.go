package scheduler

import (
	"Weather-Forecast-API/internal/repository"
	"Weather-Forecast-API/internal/services/notification"
	"Weather-Forecast-API/internal/weatherprovider"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	notifService    notification.NotificationService
	weatherProvider weatherprovider.WeatherProvider
	clock           Clock
	requestTimeout  time.Duration
}

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (r realClock) Now() time.Time {
	return time.Now()
}

func NewScheduler(
	notifService notification.NotificationService,
	weatherProvider weatherprovider.WeatherProvider,
	requestTimeout time.Duration,
) *Scheduler {
	return &Scheduler{
		notifService:    notifService,
		weatherProvider: weatherProvider,
		clock:           realClock{},
		requestTimeout:  requestTimeout,
	}
}

func (s *Scheduler) Start() (*cron.Cron, error) {
	log.Println("[Scheduler] Starting scheduler...")

	template, err := repository.GetTemplateByName("weather_update")
	if err != nil {
		return nil, fmt.Errorf("[Scheduler] Failed to get template: %v", err)
	}

	cronScheduler := cron.New()

	_, err = cronScheduler.AddFunc("@every 1m", func() {
		ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
		defer cancel()

		log.Println("[Scheduler] Checking subscriptions...")

		subs := repository.GetDueSubscriptions()
		log.Printf("[Scheduler] Found %d due subscriptions\n", len(subs))

		for _, sub := range subs {
			log.Printf("[Scheduler] Processing subscription %d for city %s\n", sub.ID, sub.City)

			weatherData, err := s.weatherProvider.GetWeatherByCity(ctx, sub.City)
			if err != nil {
				log.Printf("[Scheduler] Error fetching weather for %s: %v\n", sub.City, err)
				continue
			}
			log.Printf("[Scheduler] Weather data received for %s: %.1f°C, %d%% humidity\n",
				sub.City, weatherData.Temperature, int(weatherData.Humidity))

			message := template.Message
			message = strings.ReplaceAll(message, "{{ city }}", sub.City)
			message = strings.ReplaceAll(message, "{{ description }}", weatherData.Description)
			message = strings.ReplaceAll(message, "{{ temperature }}", fmt.Sprintf("%.1f", weatherData.Temperature))
			message = strings.ReplaceAll(message, "{{ humidity }}", strconv.Itoa(int(weatherData.Humidity)))

			subject := template.Subject
			subject = strings.ReplaceAll(subject, "{{ city }}", sub.City)

			log.Printf("[Scheduler] Sending notification to %s via %s\n", sub.ChannelValue, sub.ChannelType)
			err = s.notifService.SendMessage(sub.ChannelType, sub.ChannelValue, message, subject)
			if err != nil {
				log.Printf("[Scheduler] Error sending notification for subscription %d: %v\n", sub.ID, err)
				continue
			}

			nextNotification := s.clock.Now().Add(time.Duration(sub.FrequencyMinutes) * time.Minute)
			log.Printf("[Scheduler] Updating next notification time for subscription %d to %v\n", sub.ID, nextNotification)
			err = repository.UpdateNextNotification(sub.ID, nextNotification)
			if err != nil {
				log.Printf("[Scheduler] Error updating next notification time: %v\n", err)
				continue
			}
		}
	})
	if err != nil {
		return nil, fmt.Errorf("[Scheduler] Failed to add cron job: %v", err)
	}

	cronScheduler.Start()
	log.Println("[Scheduler] Started successfully")

	return cronScheduler, nil
}
