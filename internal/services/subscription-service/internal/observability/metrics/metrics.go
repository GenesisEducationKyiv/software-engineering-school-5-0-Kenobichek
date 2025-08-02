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

	err := prometheus.Register(ActiveSubscriptions)
	if err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			return fmt.Errorf("failed to register ActiveSubscriptions: %w", err)
		}
	}

	err = prometheus.Register(SubscriptionsCreated)
	if err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			return fmt.Errorf("failed to register SubscriptionsCreated: %w", err)
		}
	}

	err = prometheus.Register(SubscriptionCreationErrors)
	if err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			return fmt.Errorf("failed to register SubscriptionCreationErrors: %w", err)
		}
	}

	if !registered {
		if err := prometheus.Register(collectors.NewGoCollector()); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				return fmt.Errorf("failed to register Go collector: %w", err)
			}
		}

		if err := prometheus.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				return fmt.Errorf("failed to register process collector: %w", err)
			}
		}
		
		registered = true
	}

	return nil
}
