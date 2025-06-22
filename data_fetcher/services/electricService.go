package services

import (
	"os"
	"time"
	"zendo/data_fetcher/errors"
	"zendo/data_fetcher/model"
	"zendo/data_fetcher/utils"

	"go.uber.org/zap"
)

type IElectricService interface {
	GetDataSince(date *time.Time) (*model.LatestEnergeyResponse, error)
	Get24HrsOfData() (*[]model.LatestEnergeyResponse, error)
}

type ElectricitymapService struct {
	Http utils.IHttpClient
}

const (
	latestElectricEndpoint     string = "https://api.electricitymap.org/v3/power-breakdown/latest?zone=GB&disableEstimations=true"
	historicalElectricEndpoint string = "https://api.electricitymap.org/v3/power-breakdown/history?zone=GB&disableEstimations=true" // NOTE: This always returns 24 hrs
)

func (s *ElectricitymapService) GetDataSince(date *time.Time) (*model.LatestEnergeyResponse, error) {
	var body model.LatestEnergeyResponse
	result, err := s.Http.Get(latestElectricEndpoint, &body, &utils.HttpOptions{
		Headers: &map[string]string{
			"auth-token": os.Getenv("ELECTRICITY_MAPS_API_KEY"),
		},
	})

	if err != nil {
		zap.L().DPanic("Failed to get latest energy usage", zap.Error(err))
		return nil, &errors.HttpError{
			StatusCode: result.StatusCode,
		}
	}

	if date != nil && (date.UTC().After(body.SourceTime.UTC()) || date.UTC().Equal(body.SourceTime.UTC())) {
		// no updates
		return nil, nil
	}

	// set time for couchdb
	body.Timestamp = &body.SourceTime

	return &body, nil
}

type HistoricalPowerResponse struct {
	History []model.LatestEnergeyResponse `json:"history"`
}

func (s *ElectricitymapService) Get24HrsOfData() (*[]model.LatestEnergeyResponse, error) {
	var body HistoricalPowerResponse
	result, err := s.Http.Get(historicalElectricEndpoint, &body, &utils.HttpOptions{
		Headers: &map[string]string{
			"auth-token": os.Getenv("ELECTRICITY_MAPS_API_KEY"),
		},
	})

	if err != nil {
		zap.L().DPanic("Failed to get historical energy usage", zap.Error(err))
		return nil, &errors.HttpError{
			StatusCode: result.StatusCode,
		}
	}

	// set timestamps
	for i := range body.History {
		body.History[i].Timestamp = &body.History[i].SourceTime
	}

	return &body.History, nil
}
