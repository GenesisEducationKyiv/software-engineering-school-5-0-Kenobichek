package handlers

import (
	"net/http"
	"os"

	"Weather-Forecast-API/internal/utilities"
	"Weather-Forecast-API/internal/weather"
)

func GetWeather(writer http.ResponseWriter, request *http.Request) {
	city := request.URL.Query().Get("city")
	if city == "" {
		utilities.RespondJSON(writer, http.StatusNotFound, "City not found")

		return
	}

	provider := weather.OpenWeather{APIKey: os.Getenv("OPENWETHERMAP_API_KEY")}

	data, err := provider.GetWeather(city)
	if err != nil {
		if err.Error() == "city not found" {
			utilities.RespondJSON(writer, http.StatusNotFound, "City not found")
		} else {
			utilities.RespondJSON(writer, http.StatusBadRequest, "Failed to get weather: "+err.Error())
		}

		return
	}

	utilities.RespondDataJSON(writer, http.StatusOK, data)
}
