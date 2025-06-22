package services

import (
	"time"
	"zendo/lib_zendo/errors"
	"zendo/lib_zendo/model"
	"zendo/lib_zendo/utils"

	"go.uber.org/zap"
)

type IWeatherService interface {
	GetDataSince(date *time.Time) (*model.WeatherResponse, error)
	Get24HrsOfData() (*[]model.WeatherResponse, error)
}

type OpenMeteoWeatherService struct {
	Http utils.IHttpClient
}

const (
	// NOTE: Default to using York as a location
	latestWeatherEndpoint     string = "https://api.open-meteo.com/v1/forecast?latitude=53.9324727&longitude=-1.1204176&current=temperature_2m&current=direct_radiation&current=cloud_cover&current=wind_speed_10m"
	historicalWeatherEndpoint string = "https://api.open-meteo.com/v1/forecast?latitude=53.9324727&longitude=-1.1204176&hourly=temperature_2m&hourly=direct_radiation&hourly=cloud_cover&hourly=wind_speed_10m&past_days=1&forecast_days=1"
	weatherTimeLayout         string = "2006-01-02T15:04"
)

func (s *OpenMeteoWeatherService) GetDataSince(date *time.Time) (*model.WeatherResponse, error) {
	var body model.WeatherResponse
	result, err := s.Http.Get(latestWeatherEndpoint, &body)
	if err != nil {
		zap.L().DPanic("Failed to get latest energy usage", zap.Error(err))
		return nil, &errors.HttpError{
			StatusCode: result.StatusCode,
		}
	}

	parsedTime, err := time.Parse(weatherTimeLayout, *body.WeatherData.TimeString)
	if err != nil {
		zap.L().DPanic("Failed to parse weather time", zap.Error(err))
		return nil, &errors.DatabaseError{}
	}

	body.Timestamp = &parsedTime
	body.WeatherData.TimeString = nil

	if date != nil && (date.UTC().After(parsedTime.UTC()) || date.UTC().Equal(parsedTime.UTC())) {
		// no updates
		return nil, nil
	}

	return &body, nil
}

type HourlyData struct {
	TimeStrings          []string  `json:"time"`
	Temperatures         []float32 `json:"temperature_2m"`
	Radiation            []float32 `json:"direct_radiation"`
	CloudCoverPercentage []byte    `json:"cloud_cover"`
	WindSpeed            []float32 `json:"wind_speed_10m"`
}

type HistoricalWeatherResponse struct {
	HourlyData HourlyData `json:"hourly"`
}

func (s *OpenMeteoWeatherService) Get24HrsOfData() (*[]model.WeatherResponse, error) {
	// NOTE: The weather api is a bit clunky and filtering by time is a pain
	// Need to filter through and remove weather in the future

	var body HistoricalWeatherResponse
	result, err := s.Http.Get(historicalWeatherEndpoint, &body)
	if err != nil {
		zap.L().DPanic("Failed to get historical weather data", zap.Error(err))
		return nil, &errors.HttpError{
			StatusCode: result.StatusCode,
		}
	}

	now := time.Now().UTC()

	dataPoints := []model.WeatherResponse{}
	for i := range body.HourlyData.TimeStrings {
		// parse time
		parsedTime, err := time.Parse(weatherTimeLayout, body.HourlyData.TimeStrings[i])
		if err != nil {
			zap.L().DPanic("Failed to parse weather time", zap.Error(err))
			return nil, &errors.DatabaseError{}
		}

		if parsedTime.UTC().After(now) {
			// open meteo start in the past and go to the future so we can break here
			break
		}

		point := model.WeatherResponse{
			BaseDocument: model.BaseDocument{
				Timestamp: &parsedTime,
			},
			WeatherData: model.WeatherData{
				Temperature:       body.HourlyData.Temperatures[i],
				DirectRadiation:   body.HourlyData.Radiation[i],
				CloudCoverPercent: body.HourlyData.CloudCoverPercentage[i],
				WindSpeedKmPHr:    body.HourlyData.WindSpeed[i],
			},
		}
		dataPoints = append(dataPoints, point)
	}

	return &dataPoints, nil
}
