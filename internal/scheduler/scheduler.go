package scheduler

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"Weather-Forecast-API/internal/notifier"
	"Weather-Forecast-API/internal/repository"
	"Weather-Forecast-API/internal/weather"
	"github.com/robfig/cron/v3"
)

func StartScheduler() {
	template, err := repository.GetTemplateByName("weather_update")
	if err != nil {
		return
	}

	cron := cron.New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = cron.AddFunc("@every 1m", func() {
		log.Println("[Scheduler] Checking subscriptions...")

		subs := repository.GetDueSubscriptions()

		for _, sub := range subs {
			provider := weather.OpenWeather{APIKey: os.Getenv("OPENWEATHERMAP_API_KEY")}

			weatherData, err := provider.GetWeather(ctx, sub.City)
			if err != nil {
				log.Printf("[Scheduler] Error fetching weather for %s: %v\n", sub.City, err)
				continue
			}

			message := template.Message
			message = strings.ReplaceAll(message, "{{ city }}", sub.City)
			message = strings.ReplaceAll(message, "{{ description }}", weatherData.Description)
			message = strings.ReplaceAll(message, "{{ temperature }}", fmt.Sprintf("%.1f", weatherData.Temperature))
			message = strings.ReplaceAll(message, "{{ humidity }}", strconv.Itoa(int(weatherData.Humidity)))

			subject := template.Subject
			subject = strings.ReplaceAll(subject, "{{ city }}", sub.City)

			emailNotifier := notifier.EmailNotifier{}
			_ = emailNotifier.Send(sub.ChannelValue, message, subject)

			err = repository.UpdateNextNotification(sub.ID, time.Now().Add(time.Duration(sub.FrequencyMinutes)*time.Minute))
			if err != nil {
				log.Printf("[Scheduler] Error updating next notification for subscription %d: %v\n", sub.ID, err)
				continue
			}
		}
	})
	if err != nil {
		return
	}

	cron.Start()
	log.Println("[Scheduler] Started")
}
