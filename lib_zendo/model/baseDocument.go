package model

import "time"

const (
	ENERGY_TYPE  string = "ENERGY_DATA"
	WEATHER_TYPE string = "WEATHER_DATA"
)

type BaseDocument struct {
	Type           *string    `json:"type,omitempty"`
	Timestamp      *time.Time `json:"timestamp,omitempty"`
	HistoricalSeed bool       `json:"historicalSeed"` // This is set so that the data processing service knows to ignore it
}
