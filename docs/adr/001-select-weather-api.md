# ADR-001: Selection of Weather API for Weather-Subscription Service

## Status
- **Accepted** — June 7, 2025  
- **Author** — Backend Developer (*Kenobichek*)

---

## 1. Context

We're building a **Weather-Subscription API** that allows users to subscribe to weather updates for their city. A core dependency is choosing the right external weather API provider.

---

## 2. Options Considered

We evaluated 5 options:

| Provider         | Free Tier                 | Historical Data | Forecast Range      | Alerts | AQI  | Mapping Support | Global Coverage | SLA                     |
|------------------|---------------------------|------------------|----------------------|--------|------|------------------|------------------|--------------------------|
| OpenWeatherMap   | 1000/day, 1M/month        | 46 years         | 1min–1.5 years       | Yes    | Yes  | 15 layers         | 200k+ cities     | 99.5% (free), 99.9% (paid) |
| WeatherAPI.com   | 1M/month                  | Since 2010       | Up to 14 days        | Yes    | Yes  | —                 | Global           | 95.5%                    |
| Weatherbit.io    | 50/day                    | 25 years         | 240hr/16 days        | Yes    | Yes  | Beta             | 120k stations    | ≥95%                     |
| AccuWeather      | 50/day                    | Paid add-on      | 120hr/15 days        | Paid   | Add-on | Limited         | Global           | Contract-based           |
| Weatherstack     | 100/month                 | No               | 7 days               | No     | No   | —                 | Global (low detail) | 99%                    |

---

## 3. Decision

✅ **Choose OpenWeatherMap One Call 3.0**, because:

- **Cost-effective**: Generous free tier for development/testing, pay-as-you-go pricing for production.
- **Rich data**: Includes minute/hourly/daily forecasts, historical data, air quality index (AQI), and weather alerts.
- **Go support**: Mature Go client available (`github.com/kyroy/weather`).
- **Reliable**: SLA suitable for our operational needs.

---

## 4. Consequences

We will:

- Create a wrapper `openweathermap.Client` in `internal/adapter/weather/`.
- Store the `OPENWETHERMAP_API_KEY` in `.env`.
---

## 5. Revisit Conditions

This ADR lives in `docs/adr/001-select-weather-api.md` and should be revisited upon significant changes to API usage, cost, or business needs.
