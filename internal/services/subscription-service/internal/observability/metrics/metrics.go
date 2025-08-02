package metrics

import "github.com/prometheus/client_golang/prometheus"

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

func Register() {
    prometheus.MustRegister(
        ActiveSubscriptions,
        SubscriptionsCreated,
        SubscriptionCreationErrors,
    )
}
