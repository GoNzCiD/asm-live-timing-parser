package helpers

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

type Hotlap struct {
	Car                string
	CarId              int
	Gap                time.Duration
	Laps               int
	LapTime            time.Duration
	Max_speed          float64
	Name               string
	NumLaps            int
	PlayerId           string
	Position           int
	Tyre               string
	Sectors            Splits
	SkinPreviewPattern *string
}

type hotlapMarshall struct {
	Car                 string `json:"car"`
	CarId               int    `json:"carid"`
	Gap                 string `json:"gap"`
	Laps                int    `json:"laps"`
	LapTime             string `json:"laptime"`
	Max_speed           string `json:"max_speed"`
	Name                string `json:"name"`
	NumLaps             int    `json:"numlaps"`
	PlayerId            string `json:"playerid"`
	Position            int    `json:"position"`
	Tyre                string `json:"tyre"`
	Sector1             string `json:"sector1"`
	Sector1Best         bool   `json:"sector1_best"`
	Sector1BestAbsolute bool   `json:"sector1_best_absolute"`
	Sector2             string `json:"sector2"`
	Sector2Best         bool   `json:"sector2_best"`
	Sector2BestAbsolute bool   `json:"sector2_best_absolute"`
	Sector3             string `json:"sector3"`
	Sector3Best         bool   `json:"sector3_best"`
	Sector3BestAbsolute bool   `json:"sector3_best_absolute"`
	SkinPreviewPath     string `json:"skin_preview"`
}

func (h *Hotlap) MarshalJSON() ([]byte, error) {
	lapTimeTotalMinutes := int(math.Trunc(h.LapTime.Minutes()))
	lapTimeTotalSeconds := h.LapTime.Seconds()
	lapTimeSeconds := lapTimeTotalSeconds - (float64(lapTimeTotalMinutes) * 60)
	lapTimeStr := fmt.Sprintf("%d:%06.3f", lapTimeTotalMinutes, lapTimeSeconds)
	result := &hotlapMarshall{
		Car:                 h.Car,
		CarId:               h.CarId,
		Gap:                 fmt.Sprintf("%.3f", h.Gap.Seconds()),
		Laps:                h.Laps,
		LapTime:             lapTimeStr,
		Max_speed:           fmt.Sprintf("%.1f Km/h", h.Max_speed),
		Name:                h.Name,
		NumLaps:             h.NumLaps,
		PlayerId:            h.PlayerId,
		Position:            h.Position,
		Tyre:                h.Tyre,
		Sector1:             fmt.Sprintf("%.3f", h.Sectors.S1.SplitTime.Seconds()),
		Sector2:             fmt.Sprintf("%.3f", h.Sectors.S2.SplitTime.Seconds()),
		Sector3:             fmt.Sprintf("%.3f", h.Sectors.S3.SplitTime.Seconds()),
		Sector1Best:         h.Sectors.S1.IsDriversBest,
		Sector2Best:         h.Sectors.S2.IsDriversBest,
		Sector3Best:         h.Sectors.S3.IsDriversBest,
		Sector1BestAbsolute: h.Sectors.S1.IsBest,
		Sector2BestAbsolute: h.Sectors.S2.IsBest,
		Sector3BestAbsolute: h.Sectors.S3.IsBest,
	}

	if h.SkinPreviewPattern != nil {
		result.SkinPreviewPath = fmt.Sprintf(*h.SkinPreviewPattern, h.PlayerId)
	}

	return json.Marshal(&result)
}
