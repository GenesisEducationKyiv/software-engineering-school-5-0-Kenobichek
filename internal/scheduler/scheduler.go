package scheduler

import (
	"Weather-Forecast-API/internal/handlers/weather"
	"Weather-Forecast-API/internal/repository/subscriptions"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type notificationManager interface {
	SendWeatherUpdate(channel string, recipient string, metrics weather.Metrics) error
}

type clockManager interface {
	Now() time.Time
}

type subscriptionManager interface {
	GetSubscriptionByToken(token string) (*subscriptions.Info, error)
	GetDueSubscriptions() []subscriptions.Info
	UpdateNextNotification(id int, next time.Time) error
}

type Scheduler struct {
	notifService     notificationManager
	subService       subscriptionManager
	weatherProvChain weatherChainHandler
	clock            clockManager
	requestTimeout   time.Duration
	cronScheduler    *cron.Cron
	mu               sync.Mutex
}

type realClock struct{}

func (r realClock) Now() time.Time {
	return time.Now()
}

type weatherChainHandler interface {
	GetWeatherByCity(ctx context.Context, city string) (weather.Metrics, error)
}

func NewScheduler(
	notifService notificationManager,
	subService subscriptionManager,
	weatherProvChain weatherChainHandler,
	requestTimeout time.Duration,
) *Scheduler {
	return &Scheduler{
		notifService:     notifService,
		subService:       subService,
		weatherProvChain: weatherProvChain,
		clock:            realClock{},
		requestTimeout:   requestTimeout,
	}
}

func (s *Scheduler) Start() (*cron.Cron, error) {
	log.Println("[Scheduler] Starting scheduler...")

	s.mu.Lock()
	defer s.mu.Unlock()

	s.cronScheduler = cron.New()

	_, err := s.cronScheduler.AddFunc("@every 1m", func() {
		ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
		defer cancel()

		log.Println("[Scheduler] Checking subscriptions...")

		subs := s.subService.GetDueSubscriptions()
		log.Printf("[Scheduler] Found %d due subscriptions\n", len(subs))

		for _, sub := range subs {
			log.Printf("[Scheduler] Processing subscription %d for city %s\n", sub.ID, sub.City)

			weatherData, err := s.weatherProvChain.GetWeatherByCity(ctx, sub.City)
			if err != nil {
				log.Printf("[Scheduler] Error fetching weather for %s: %v\n", sub.City, err)
				continue
			}
			log.Printf("[Scheduler] Weather data received for %s: %.1f°C, %d%% humidity\n",
				sub.City, weatherData.Temperature, int(weatherData.Humidity))

			log.Printf("[Scheduler] Sending notification to %s via %s\n", sub.ChannelValue, sub.ChannelType)
			err = s.notifService.SendWeatherUpdate(sub.ChannelType, sub.ChannelValue, weatherData)
			if err != nil {
				log.Printf("[Scheduler] Error sending notification for subscription %d: %v\n", sub.ID, err)
				continue
			}

			nextNotification := s.clock.Now().Add(time.Duration(sub.FrequencyMinutes) * time.Minute)
			log.Printf("[Scheduler] Updating next notification time for subscription %d to %v\n", sub.ID, nextNotification)
			err = s.subService.UpdateNextNotification(sub.ID, nextNotification)
			if err != nil {
				log.Printf("[Scheduler] Error updating next notification time: %v\n", err)
				continue
			}
		}
	})
	if err != nil {
		return nil, fmt.Errorf("[Scheduler] Failed to add cron job: %v", err)
	}

	s.cronScheduler.Start()
	log.Println("[Scheduler] Started successfully")

	return s.cronScheduler, nil
}

func (s *Scheduler) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cronScheduler != nil {
		log.Println("[Scheduler] Stopping scheduler...")
		ctx := s.cronScheduler.Stop()
		<-ctx.Done()
		log.Println("[Scheduler] Stopped successfully")
	}

	return nil
}
