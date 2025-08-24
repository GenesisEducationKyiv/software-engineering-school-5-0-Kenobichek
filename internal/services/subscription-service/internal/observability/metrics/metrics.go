package metrics

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	ActiveSubscriptions = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "subscription_service",
		Name:      "active_subscriptions",
		Help:      "Current number of active subscriptions.",
	})

	SubscriptionsCreated = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "subscription_service",
		Name:      "subscriptions_created_total",
		Help:      "Total number of subscriptions successfully created.",
	})

	SubscriptionCreationErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "subscription_service",
		Name:      "subscription_creation_errors_total",
		Help:      "Total number of errors that occurred while creating subscriptions.",
	})
)

var registered bool
var registerMutex sync.Mutex

func Register() error {
	registerMutex.Lock()
	defer registerMutex.Unlock()

	if registered {
		return nil
	}

	metrics := []prometheus.Collector{
		ActiveSubscriptions,
		SubscriptionsCreated,
		SubscriptionCreationErrors,
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	}
		
	for _, metric := range metrics {
		if err := prometheus.Register(metric); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				return fmt.Errorf("failed to register metric %T: %w", metric, err)
			}
		}
	}
	
	registered = true
	
	return nil
}
