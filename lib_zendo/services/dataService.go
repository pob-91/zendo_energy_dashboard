package services

import (
	"os"
	"strings"
	"time"
	"zendo/lib_zendo/errors"
	"zendo/lib_zendo/model"
	"zendo/lib_zendo/utils"

	"go.uber.org/zap"
)

type IDataService interface {
	PostLatestData(energy *model.LatestEnergeyResponse, weather *model.WeatherResponse) error
	SeedHistoricalData(energyData *[]model.LatestEnergeyResponse, weatherData *[]model.WeatherResponse) error
	GetLatestWeatherDate() (*time.Time, error)
	GetLatestEnergyDate() (*time.Time, error)
	GetLatestMetric() (*model.Metric, error)
	Get24HoursOfMetrics() (*[]model.Metric, error)
}

type CouchDBDataService struct {
	Http utils.IHttpClient
}

func (s *CouchDBDataService) PostLatestData(energy *model.LatestEnergeyResponse, weather *model.WeatherResponse) error {
	// NOTE: This may need to become more complex but for now we are just going to add a type property and post

	if energy == nil && weather == nil {
		zap.L().Panic("Neither energy update nor weather update was set")
		return &errors.DatabaseError{}
	}

	payloadSlice := []any{}
	if energy != nil {
		energy.BaseDocument.Type = utils.StringPointer(model.ENERGY_TYPE)
		payloadSlice = append(payloadSlice, energy)
	}
	if weather != nil {
		weather.BaseDocument.Type = utils.StringPointer(model.WEATHER_TYPE)
		payloadSlice = append(payloadSlice, weather)
	}

	payload := map[string][]any{
		"docs": payloadSlice,
	}
	if _, err := s.Http.Post(bulkDocsUrl(), payload, nil); err != nil {
		zap.L().DPanic("Failed to post latest data", zap.Error(err))
		return &errors.DatabaseError{}
	}

	return nil
}

func (s *CouchDBDataService) SeedHistoricalData(energyData *[]model.LatestEnergeyResponse, weatherData *[]model.WeatherResponse) error {
	if energyData == nil || weatherData == nil {
		zap.L().Panic("Neither energy data nor weather data was set")
		return &errors.DatabaseError{}
	}

	// set types and historical seed flag
	energyLen := len(*energyData)
	docs := make([]any, energyLen+len(*weatherData))
	for i, x := range *energyData {
		x.Type = utils.StringPointer(model.ENERGY_TYPE)
		x.HistoricalSeed = true
		docs[i] = x
	}
	for i, x := range *weatherData {
		x.Type = utils.StringPointer(model.WEATHER_TYPE)
		x.HistoricalSeed = true
		docs[i+energyLen] = x
	}

	payload := map[string][]any{
		"docs": docs,
	}
	if _, err := s.Http.Post(bulkDocsUrl(), payload, nil); err != nil {
		zap.L().DPanic("Failed to post seed data", zap.Error(err))
		return &errors.DatabaseError{}
	}

	return nil
}

type CouchDBViewDocMetadata struct {
	Id  string             `json:"id"`
	Doc model.BaseDocument `json:"doc"`
}

type CouchDBViewResponse struct {
	Rows []CouchDBViewDocMetadata `json:"rows"`
}

func (s *CouchDBDataService) GetLatestWeatherDate() (*time.Time, error) {
	var body CouchDBViewResponse
	result, err := s.Http.Get(latestWeatherUrl(), &body)
	if err != nil {
		zap.L().DPanic("Failed to get latest weather date", zap.Error(err))
		return nil, &errors.HttpError{
			StatusCode: result.StatusCode,
		}
	}

	if len(body.Rows) == 0 {
		// no data
		return nil, nil
	}

	return body.Rows[0].Doc.Timestamp, nil
}

func (s *CouchDBDataService) GetLatestEnergyDate() (*time.Time, error) {
	var body CouchDBViewResponse
	result, err := s.Http.Get(latestEnergyUrl(), &body)
	if err != nil {
		zap.L().DPanic("Failed to get latest weather date", zap.Error(err))
		return nil, &errors.HttpError{
			StatusCode: result.StatusCode,
		}
	}

	if len(body.Rows) == 0 {
		// no data
		return nil, nil
	}

	return body.Rows[0].Doc.Timestamp, nil
}

// TODO: This is not the cleanest, want better inheritance
type CouchDBViewMetricMetadata struct {
	Id  string       `json:"id"`
	Doc model.Metric `json:"doc"`
}

type CouchDBMetricViewResponse struct {
	Rows []CouchDBViewMetricMetadata `json:"rows"`
}

func (s *CouchDBDataService) GetLatestMetric() (*model.Metric, error) {
	var body CouchDBMetricViewResponse
	result, err := s.Http.Get(latestMetricUrl(), &body)
	if err != nil {
		zap.L().DPanic("Failed to get latest metric", zap.Error(err))
		return nil, &errors.HttpError{
			StatusCode: result.StatusCode,
		}
	}

	if len(body.Rows) == 0 {
		// no data
		return nil, nil
	}

	return &body.Rows[0].Doc, nil
}

func (s *CouchDBDataService) Get24HoursOfMetrics() (*[]model.Metric, error) {
	var body CouchDBMetricViewResponse
	result, err := s.Http.Get(last24HoursMetricUrl(), &body)
	if err != nil {
		zap.L().DPanic("Failed to get 24 hours of metrics", zap.Error(err))
		return nil, &errors.HttpError{
			StatusCode: result.StatusCode,
		}
	}

	if len(body.Rows) == 0 {
		// no data
		return &[]model.Metric{}, nil
	}

	data := make([]model.Metric, len(body.Rows))
	for i := range body.Rows {
		data[i] = body.Rows[i].Doc
	}

	return &data, nil
}

// private

func baseUrl() string {
	var builder strings.Builder
	builder.WriteString("http://")
	builder.WriteString(os.Getenv("COUCHDB_USER"))
	builder.WriteString(":")
	builder.WriteString(os.Getenv("COUCHDB_PASSWORD"))
	builder.WriteString("@")
	builder.WriteString(os.Getenv("COUCHDB_URL"))
	builder.WriteString("/")
	builder.WriteString(os.Getenv("COUCHDB_DB"))
	return builder.String()
}

func bulkDocsUrl() string {
	var builder strings.Builder
	builder.WriteString(baseUrl())
	builder.WriteString("/_bulk_docs")
	return builder.String()
}

func latestWeatherUrl() string {
	var builder strings.Builder
	builder.WriteString(baseUrl())
	builder.WriteString("/_design/views/_view/weather_by_time?include_docs=true&descending=true&limit=1")
	return builder.String()
}

func latestMetricUrl() string {
	var builder strings.Builder
	builder.WriteString(baseUrl())
	builder.WriteString("/_design/views/_view/aggregated_by_time?include_docs=true&descending=true&limit=1")
	return builder.String()
}

func last24HoursMetricUrl() string {
	yesterday := time.Now().Add(-time.Hour * 24)

	var builder strings.Builder
	builder.WriteString(baseUrl())
	builder.WriteString("/_design/views/_view/aggregated_by_time?include_docs=true&descending=true&startKey=\"")
	builder.WriteString(yesterday.Format(time.RFC3339))
	builder.WriteString("\"")
	return builder.String()
}

func latestEnergyUrl() string {
	var builder strings.Builder
	builder.WriteString(baseUrl())
	builder.WriteString("/_design/views/_view/aggregated_by_time?include_docs=true&descending=true&limit=1")
	return builder.String()
}
