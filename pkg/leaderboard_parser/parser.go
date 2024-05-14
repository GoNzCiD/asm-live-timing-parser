package leaderboard_parser

import (
	"acsm-live_timing-parser/pkg/helpers"
	"encoding/json"
	"sort"
	"time"
)

func ReadJson(jsonStr string) (*helpers.LeaderBoard, error) {
	result := helpers.LeaderBoard{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	return &result, err
}

func ExtractHotlaps(drivers []helpers.Driver) []helpers.Hotlap {
	var result []helpers.Hotlap
	for _, driver := range drivers {
		for _, car := range driver.Cars {
			lap := helpers.Hotlap{
				Car:       car.CarName,
				CarId:     driver.CarInfo.CarID,
				Laps:      car.NumLaps,
				LapTime:   time.Duration(car.BestLap),
				Max_speed: car.TopSpeedBestLap,
				Name:      driver.CarInfo.DriverName,
				NumLaps:   driver.TotalNumLaps,
				PlayerId:  driver.CarInfo.DriverGUID,
				Tyre:      car.TyreBestLap,
				Sectors:   car.BestLapSplits,
			}
			result = append(result, lap)
		}
	}

	return result
}

func SortAndCalculateData(hotlaps []helpers.Hotlap, skinPreviewPattern *string) []helpers.Hotlap {
	sort.Slice(hotlaps, func(i, j int) bool {
		if hotlaps[i].LapTime == 0 {
			return false
		}
		return hotlaps[i].LapTime < hotlaps[j].LapTime
	})

	if len(hotlaps) > 0 {
		best := hotlaps[0].LapTime
		for i := range hotlaps {
			hotlaps[i].Position = i + 1
			gap := hotlaps[i].LapTime - best
			if gap < 0 {
				gap = 0
			}
			hotlaps[i].Gap = gap
			hotlaps[i].SkinPreviewPattern = skinPreviewPattern
		}
	}

	return hotlaps
}
