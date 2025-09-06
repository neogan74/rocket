package models

import "time"

type Weather struct {
	City        float64   `json:"city"`
	Temperature float64   `json:"temperature"`
	Humidity    int       `json:"humidity"`
	UpdatedAt   time.Time `json:"updated_at"`
}
