package scheduler

import (
	"Weather-Forecast-API/internal/notifier"
	"Weather-Forecast-API/internal/repository"
	"Weather-Forecast-API/internal/weather"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

func StartScheduler() {
	template, err := repository.GetTemplateByName("weather_update")
	if err != nil {
		return
	}

	c := cron.New()

	c.AddFunc("@every 1m", func() {
		log.Println("[Scheduler] Checking subscriptions...")

		subs := repository.GetDueSubscriptions()

		for _, sub := range subs {
			provider := weather.OpenWeather{APIKey: os.Getenv("OPENWETHERMAP_API_KEY")}
			weather_data, err := provider.GetWeather(sub.City)

			if err != nil {
				log.Printf("[Scheduler] Error fetching weather for %s: %v\n", sub.City, err)
				continue
			}

			message := template.Message
			message = strings.ReplaceAll(message, "{{ city }}", sub.City)
			message = strings.ReplaceAll(message, "{{ description }}", weather_data.Description)
			message = strings.ReplaceAll(message, "{{ temperature }}", fmt.Sprintf("%.1f", weather_data.Temperature))
			message = strings.ReplaceAll(message, "{{ humidity }}", fmt.Sprintf("%d", int(weather_data.Humidity)))

			subject := template.Subject
			subject = strings.ReplaceAll(subject, "{{ city }}", sub.City)

			notifier := notifier.EmailNotifier{}
			_ = notifier.Send(sub.ChannelValue, message, subject)

			repository.UpdateNextNotification(sub.ID, time.Now().Add(time.Duration(sub.FrequencyMinutes)*time.Minute))
		}
	})

	c.Start()
	log.Println("[Scheduler] Started")
}
