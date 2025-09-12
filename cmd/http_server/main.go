package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	httpPort = ":8080"
)

func main() {
	r := chi.NewRouter()

	r.Get("/weather", getWeatherHandler)
}

func getWeatherHandler(w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		city := chi.URLParam(r, "city")
		if city == "" {
			http.Error(w, "City parameter is required", http.StatusBadRequest)
			return
		}

		weather := storage.GetWeather(city)
		if weather == nil {
			http.Error(w, "Weather data not found", http.StatusNotFound)
			return
		}

		render.JSON(w, r, weather)
	}
}
