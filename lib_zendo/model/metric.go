package model

type CorrelationData struct {
	SolarIrradianceVsSolarProductionCorrelation float32 `json:"solar_irradiance_vs_solar_production_correlation"`
	TemperatureVsConsumptionCorrelation         float32 `json:"temperature_vs_consumption_correlation"`
}

type Metric struct {
	BaseDocument

	TotalProduction      uint32                    `json:"totalProduction"`
	TotalConsumption     uint32                    `json:"totalConsumption"`
	NetBalance           int32                     `json:"netBalance"`
	WeatherData          WeatherData               `json:"weatherData"`
	PowerProductionData  PowerProductionBreakdown  `json:"powerProductionData"`
	PowerConsumptionData PowerConsumptionBreakdown `json:"powerConsumptionData"`
	CorrelationData      CorrelationData           `json:"correlations"`
}
