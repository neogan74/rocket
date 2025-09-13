package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/neogan74/rocket/pkg/models"
)

const (
	httpPort        = ":8080"
	shutdownTimeout = 5 * time.Second
)

func main() {
	storage := models.NewInMemoryWeatherStorage()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api/v1/weather", func(r chi.Router) {
		r.Get("/{city}", getWeatherHandler(storage))
		r.Put("/{city}", updateWeatherHandler(storage))
	})

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: 10,
	}

	go func() {
		log.Printf(" HTTP server listening on %s", httpPort)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start HTTP server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")

}

func getWeatherHandler(storage *models.InMemoryWeatherStorage) http.HandlerFunc {
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

func updateWeatherHandler(storage *models.InMemoryWeatherStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		city := chi.URLParam(r, "city")
		if city == "" {
			http.Error(w, "City parameter is required", http.StatusBadRequest)
			return
		}

		var weatherUpdate models.Weather
		if err := json.NewDecoder(r.Body).Decode(&weatherUpdate); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		weatherUpdate.City = city

		weatherUpdate.UpdatedAt = time.Now()
		storage.UpdateWeather(&weatherUpdate)

		w.WriteHeader(http.StatusNoContent)
	}
}
