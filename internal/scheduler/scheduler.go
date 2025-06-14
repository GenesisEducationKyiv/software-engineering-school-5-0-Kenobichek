package scheduler

import (
	weather "Weather-Forecast-API/internal/external/openweather"
	"Weather-Forecast-API/internal/repository"
	notificationService "Weather-Forecast-API/internal/services/notification"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

func StartScheduler() {
	log.Println("[Scheduler] Starting scheduler...")

	template, err := repository.GetTemplateByName("weather_update")
	if err != nil {
		log.Printf("[Scheduler] Failed to get template: %v\n", err)
		return
	}

	cronScheduler := cron.New()

	_, err = cronScheduler.AddFunc("@every 1m", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		log.Println("[Scheduler] Checking subscriptions...")

		subs := repository.GetDueSubscriptions()
		log.Printf("[Scheduler] Found %d due subscriptions\n", len(subs))

		provider := weather.NewOpenWeatherProvider(os.Getenv("OPENWEATHERMAP_API_KEY"))
		notificationSender := notificationService.NewNotificationService()

		for _, sub := range subs {
			log.Printf("[Scheduler] Processing subscription %d for city %s\n", sub.ID, sub.City)

			weatherData, err := provider.GetWeatherByCity(ctx, sub.City)
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
			err = notificationSender.SendMessage(sub.ChannelType, sub.ChannelValue, message, subject)
			if err != nil {
				log.Printf("[Scheduler] Error sending notification for subscription %d: %v\n", sub.ID, err)
				continue
			}

			nextNotification := time.Now().Add(time.Duration(sub.FrequencyMinutes) * time.Minute)
			log.Printf("[Scheduler] Updating next notification time for subscription %d to %v\n", sub.ID, nextNotification)
			err = repository.UpdateNextNotification(sub.ID, nextNotification)
			if err != nil {
				log.Printf("[Scheduler] Error updating next notification time: %v\n", err)
				continue
			}
		}
	})
	if err != nil {
		log.Printf("[Scheduler] Failed to add cron job: %v\n", err)
		return
	}

	cronScheduler.Start()
	log.Println("[Scheduler] Started successfully")
}
