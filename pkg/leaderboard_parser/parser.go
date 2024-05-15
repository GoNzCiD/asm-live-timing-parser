package leaderboard_parser

import (
	"acsm-live_timing-parser/pkg/helpers"
	"encoding/json"
	"math"
	"sort"
	"time"
)

func ReadJson(jsonStr string) (*helpers.LeaderBoard, error) {
	result := helpers.LeaderBoard{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	return &result, err
}

func ExtractHotlaps(drivers []helpers.Driver) ([]helpers.Hotlap, helpers.Splits) {
	var bestSectors []time.Duration
	bestSectors[0] = math.MaxInt64
	bestSectors[1] = math.MaxInt64
	bestSectors[2] = math.MaxInt64

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

			if car.BestSplits.S1.SplitTime < bestSectors[0] {
				bestSectors[0] = car.BestSplits.S1.SplitTime
			}
			if car.BestSplits.S2.SplitTime < bestSectors[1] {
				bestSectors[1] = car.BestSplits.S2.SplitTime
			}
			if car.BestSplits.S3.SplitTime < bestSectors[2] {
				bestSectors[2] = car.BestSplits.S3.SplitTime
			}
		}
	}

	return result, helpers.Splits{
		S1: helpers.Split{
			SplitIndex: 0,
			SplitTime:  bestSectors[0],
		},
		S2: helpers.Split{
			SplitIndex: 1,
			SplitTime:  bestSectors[1]},
		S3: helpers.Split{
			SplitIndex: 2,
			SplitTime:  bestSectors[2]},
	}
}

func SortAndCalculateData(hotlaps []helpers.Hotlap, skinPreviewPattern *string, bestSectors helpers.Splits) []helpers.Hotlap {
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

			hotlaps[i].Sectors.S1.IsBest = hotlaps[i].Sectors.S1.SplitTime == bestSectors.S1.SplitTime
			hotlaps[i].Sectors.S2.IsBest = hotlaps[i].Sectors.S2.SplitTime == bestSectors.S2.SplitTime
			hotlaps[i].Sectors.S3.IsBest = hotlaps[i].Sectors.S3.SplitTime == bestSectors.S3.SplitTime
		}
	}

	return hotlaps
}
