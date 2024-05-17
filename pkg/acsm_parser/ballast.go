package acsm_parser

import (
	"acsm-live_timing-parser/pkg/helpers"
	"encoding/json"
)

type DriverBallast struct {
	DriverGUID string `json:"driver_guid"`
	Ballast    int    `json:"ballast"`
}

type BallastCollection []DriverBallast

func SaveBallast(results []helpers.ACSMResult, saveToPath string, ballastDistribution []int) error {
	var ballast BallastCollection
	for i, result := range results {
		if i >= len(ballastDistribution) {
			break
		}
		ballast = append(ballast, DriverBallast{
			DriverGUID: result.DriverGUID,
			Ballast:    ballastDistribution[i],
		})
	}

	ballastBytes, err := json.Marshal(ballast)
	if err != nil {
		return err
	}

	err = helpers.SaveToFile(string(ballastBytes), saveToPath)
	if err != nil {
		return err
	}

	return nil
}

func ReadAndParseBallast(path string) (map[string]int, error) {
	ballastStr, err := helpers.ReadFromFile(path)
	if err != nil {
		return nil, err
	}
	ballast := BallastCollection{}
	err = json.Unmarshal([]byte(ballastStr), &ballast)
	if err != nil {
		return nil, err
	}

	result := map[string]int{}
	for _, driverBallast := range ballast {
		result[driverBallast.DriverGUID] = driverBallast.Ballast
	}

	return result, nil
}
