package app

import (
	"fmt"
	"log"

	"internal/services/weather-service/config"
	"internal/services/weather-service/infrastructure"
	"internal/services/weather-service/internalhttpclient"
	"internal/services/weather-service/provider"
	"internal/services/weather-service/server"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	redisAddr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	cache, err := infrastructure.NewRedisCache(redisAddr, cfg.Redis.Password, cfg.Redis.DB, cfg.Redis.CacheTTL)
	if err != nil {
		return err
	}
	defer cache.Close()

	httpClient := internalhttpclient.New()

	geo := infrastructure.NewGeocodingService(httpClient, cfg.OpenWeather.GeocodingAPIURL, cfg.OpenWeather.APIKey)
	openWeather := infrastructure.NewOpenWeatherAPI(httpClient, cfg.OpenWeather.WeatherAPIURL, cfg.OpenWeather.APIKey)
	openWeatherProvider := provider.NewOpenWeatherProvider(geo, openWeather)

	weatherAPI := infrastructure.NewWeatherAPIProvider(httpClient, cfg.WeatherAPI.URL, cfg.WeatherAPI.APIKey)
	weatherAPIProvider := provider.NewWeatherAPIProvider(weatherAPI)

	openWeatherChain := provider.NewChainWeatherProvider(openWeatherProvider)
	weatherAPIChain := provider.NewChainWeatherProvider(weatherAPIProvider)
	weatherAPIChain.SetNext(openWeatherChain)

	cachedProvider := provider.NewCachedWeatherProvider(weatherAPIChain, cache)

	address := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("weather-service starting on %s", address)
	if err := server.RunGRPCServer(address, cachedProvider); err != nil {
		log.Printf("weather-service exited with error: %v", err)
		return err
	}
	log.Println("weather-service stopped")
	return nil
}
