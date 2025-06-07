# ADR‑001: Selecting a Weather‑Data API Provider for the **Weather‑Subscription API**

**Status**: Accepted – 7 June 2025  
**Decision Maker**: Backend Dev (*Kenobichek*)
---

## Context

The **Weather‑Subscription API** - a simple API that lets users subscribe to weather updates for their city  

Five commercial weather APIs were shortlisted and evaluated:

| Criterion | **OpenWeatherMap** (chosen) | WeatherAPI.com | Weatherbit.io | AccuWeather | Weatherstack |
|-----------|-----------------------------|----------------|---------------|-------------|--------------|
| Free tier & limits | 1 000 calls/day; 60 req/min; **1 M req/mo** | 1 M req/mo; 3‑day forecast | 50 req/day | 50 req/day | 100 req/mo |
| Paid price entry | Pay‑as‑you‑go **$0.0015/call** | **$7/mo** (3 M) | **$40/mo** | **$25/mo** | **$9.99/mo** |
| Historical data | 46 + years | 2010→ | 25 y | Add‑on | None |
| Forecast range | Minute‑by‑minute 1 h; hourly 48 h; daily 8 d; long‑term 1.5 y | Up to 14 d | Hourly 240 h; daily 16 d | Hourly 120 h; daily 15 d | 7 d |
| Severe‑weather alerts | Yes | Yes | Yes | Yes (paid) | No |
| Air‑quality data | Yes (global AQI) | Yes | Yes | Add‑on | No |
| Maps & tiles | 15 static/tile layers | — | Beta tiles | Limited radar | — |
| Global coverage | 200 k + cities | Global | 120 k stations | Global | Global (coarse) |
| SLA / uptime | 99.5 % free; 99.9 % paid | 95.5 % free | ≥ 95 % | Contract | 99 % |

---

## Decision

Adopt **OpenWeatherMap One Call 3.0 API.**

* **Cost‑effective scaling** – free 1 000 daily calls cover dev/testing; pay‑as‑you‑go eliminates big plan jumps.
* **Rich dataset** – 46 + years of history, 1.5‑year outlook, AQI and government alerts cover all planned features.
* **Strong Go ecosystem** – community package [`github.com/kyroy/weather`](https://pkg.go.dev/github.com/kyroy/weather) shortens integration time.
* **Mature SLAs** – 99.5 % on free, 99.9 % on paid tiers satisfy our non‑functional requirements.
---

## Consequences

* **Implementation** – add an `openweathermap.Client` wrapper in `internal/adapter/weather/`.
* **Configuration** – store `OWM_API_KEY` in `.env`; inject via Viper.
* **Monitoring** – expose Prometheus gauge `owm_remaining_quota`; alert at < 20 % daily quota.
* **Future** – if monthly cost exceeds **$500** or data quality issues arise, re‑evaluate WeatherAPI (closest feature parity).

---

*This ADR resides at `docs/adr/001-select-weather-api.md` and may be revisited as project requirements evolve.*