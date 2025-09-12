package models

import "sync"

// WeatherStorage represents the weather data for a specific city.
type WeatherStorage interface {
	GetWeather(city string) *Weather
	UpdateWeather(city string, weather *Weather)
}

// InMemoryWeatherStorage is an in-memory implementation of WeatherStorage.
type InMemoryWeatherStorage struct {
	mu       sync.RWMutex
	weathers map[string]*Weather
}

// NewInMemoryWeatherStorage creates a new instance of InMemoryWeatherStorage.
func NewInMemoryWeatherStorage() *InMemoryWeatherStorage {
	return &InMemoryWeatherStorage{
		weathers: make(map[string]*Weather),
	}
}

// GetWeather retrieves the weather data for a specific city.
func (s *InMemoryWeatherStorage) GetWeather(city string) *Weather {
	s.mu.RLock()
	defer s.mu.RUnlock()

	weather, exists := s.weathers[city]
	if !exists {
		return nil
	}
	return weather
}

// UpdateWeather updates or adds the weather data for a specific city.
func (s *InMemoryWeatherStorage) UpdateWeather(city string, weather *Weather) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.weathers[city] = weather
}
