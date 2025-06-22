package routes

import (
	"encoding/json"
	"net/http"
	"zendo/lib_zendo/services"

	"go.uber.org/zap"
)

type DataRoutes struct {
	DataService services.IDataService
}

func (r *DataRoutes) GetLatestMetric(resp http.ResponseWriter, req *http.Request) {
	metric, err := r.DataService.GetLatestMetric()
	if err != nil {
		zap.L().DPanic("Failed to get latest metric", zap.Error(err))
		resp.WriteHeader(http.StatusFailedDependency)
		return
	}

	if metric == nil {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(resp).Encode(metric); err != nil {
		zap.L().DPanic("Failed to encode metric", zap.Error(err))
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (r *DataRoutes) GetTimeSeriesMetrics(resp http.ResponseWriter, req *http.Request) {
	// TODO: Really want pagination here
	metrics, err := r.DataService.Get24HoursOfMetrics()
	if err != nil {
		zap.L().DPanic("Failed to get time serires metrics", zap.Error(err))
		resp.WriteHeader(http.StatusFailedDependency)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(resp).Encode(metrics); err != nil {
		zap.L().DPanic("Failed to encode metrics", zap.Error(err))
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
}
