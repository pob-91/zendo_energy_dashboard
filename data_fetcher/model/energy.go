package model

import "time"

type PowerProductionBreakdown struct {
	Nuclear          *uint32 `json:"nuclear,omitempty"`
	Geothermal       *uint32 `json:"geothermal,omitempty"`
	Biomass          *uint32 `json:"biomass,omitempty"`
	Coal             *uint32 `json:"coal,omitempty"`
	Wind             *uint32 `json:"wind,omitempty"`
	Solar            *uint32 `json:"solar,omitempty"`
	Hydro            *uint32 `json:"hydro,omitempty"`
	Gas              *uint32 `json:"gas,omitempty"`
	Oil              *uint32 `json:"oil,omitempty"`
	Unknown          *uint32 `json:"unknown,omitempty"`
	HydroDischarge   *uint32 `json:"hydro discharge,omitempty"`
	BatteryDischarge *uint32 `json:"battery discharge,omitempty"`
}

type PowerConsumptionBreakdown struct {
	Nuclear          *uint32 `json:"nuclear,omitempty"`
	Geothermal       *uint32 `json:"geothermal,omitempty"`
	Biomass          *uint32 `json:"biomass,omitempty"`
	Coal             *uint32 `json:"coal,omitempty"`
	Wind             *uint32 `json:"wind,omitempty"`
	Solar            *uint32 `json:"solar,omitempty"`
	Hydro            *uint32 `json:"hydro,omitempty"`
	Gas              *uint32 `json:"gas,omitempty"`
	Oil              *uint32 `json:"oil,omitempty"`
	Unknown          *uint32 `json:"unknown,omitempty"`
	HydroDischarge   *uint32 `json:"hydro discharge,omitempty"`
	BatteryDischarge *uint32 `json:"battery discharge,omitempty"`
}

type LatestEnergeyResponse struct {
	BaseDocument
	SourceTime                time.Time                 `json:"datetime"`
	PowerProductionBreakdown  PowerProductionBreakdown  `json:"powerProductionBreakdown"`
	PowerConsumptionBreakdown PowerConsumptionBreakdown `json:"powerConsumptionBreakdown"`
	PowerProductionTotal      uint32                    `json:"powerProductionTotal"`
	PowerConsumptionTotal     uint32                    `json:"powerConsumptionTotal"`
}
