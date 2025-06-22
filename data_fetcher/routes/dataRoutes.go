package routes

import (
	"net/http"
	"zendo/data_fetcher/services"

	"go.uber.org/zap"
)

type DataRoutes struct {
	ElectricService services.IElectricService
	WeatherService  services.IWeatherService
	DataService     services.IDataService
}

func (r *DataRoutes) GetLatest(resp http.ResponseWriter, req *http.Request) {
	latestWeatherUpdate, err := r.DataService.GetLatestWeatherDate()
	if err != nil {
		zap.L().DPanic("Failed to get latest weather update time")
		resp.WriteHeader(http.StatusFailedDependency)
		return
	}
	latestEnergyUpdate, err := r.DataService.GetLatestEnergyDate()
	if err != nil {
		zap.L().DPanic("Failed to get latest weather update time")
		resp.WriteHeader(http.StatusFailedDependency)
		return
	}

	latestEnergy, err := r.ElectricService.GetDataSince(latestEnergyUpdate)
	if err != nil {
		zap.L().Warn("Failed to get latest energy data, continuing anyway", zap.Error(err))
	}

	latestWeather, err := r.WeatherService.GetDataSince(latestWeatherUpdate)
	if err != nil {
		zap.L().Warn("Failed to get latest weather data, continuing anyway", zap.Error(err))
	}

	if latestEnergy == nil && latestWeather == nil {
		zap.L().Info("No data updates.")
		resp.WriteHeader(204)
		return
	}

	if err := r.DataService.PostLatestData(latestEnergy, latestWeather); err != nil {
		resp.WriteHeader(http.StatusFailedDependency)
		return
	}

	resp.WriteHeader(204)
}

func (r *DataRoutes) Seed24Hrs(resp http.ResponseWriter, req *http.Request) {
	historicalEnergy, err := r.ElectricService.Get24HrsOfData()
	if err != nil {
		zap.L().DPanic("Failed to get historical energy data", zap.Error(err))
		resp.WriteHeader(http.StatusFailedDependency)
		return
	}

	historicalWeather, err := r.WeatherService.Get24HrsOfData()
	if err != nil {
		zap.L().DPanic("Failed to get historical weather data", zap.Error(err))
		resp.WriteHeader(http.StatusFailedDependency)
		return
	}

	if err := r.DataService.SeedHistoricalData(historicalEnergy, historicalWeather); err != nil {
		zap.L().DPanic("Failed to add historical data to the database", zap.Error(err))
		resp.WriteHeader(http.StatusFailedDependency)
	}

	resp.WriteHeader(204)
}
