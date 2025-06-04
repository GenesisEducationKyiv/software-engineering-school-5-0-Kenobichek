package handlers

import (
	"Weather-Forecast-API/internal/utilities"
	"Weather-Forecast-API/internal/weather"
	"net/http"
	"os"
)

func GetWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		utilities.RespondJSON(w, http.StatusNotFound, "City not found")
		return
	}

	provider := weather.OpenWeather{APIKey: os.Getenv("OPENWETHERMAP_API_KEY")}
	data, err := provider.GetWeather(city)

	if err != nil {
		if err.Error() == "city not found" {
			utilities.RespondJSON(w, http.StatusNotFound, "City not found")
		} else {
			utilities.RespondJSON(w, http.StatusBadRequest, "Failed to get weather: "+err.Error())
		}
		return
	}

	utilities.RespondDataJSON(w, http.StatusOK, data)
}
