package model

type WeatherData struct {
	TimeString        *string `json:"time,omitempty"`
	Temperature       float32 `json:"temperature_2m"`
	DirectRadiation   float32 `json:"direct_radiation"`
	CloudCoverPercent byte    `json:"cloud_cover"`
	WindSpeedKmPHr    float32 `json:"wind_speed_10m"`
}

type WeatherResponse struct {
	BaseDocument
	WeatherData WeatherData `json:"current"`
}
