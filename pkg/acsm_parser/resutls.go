package acsm_parser

import (
	"acsm-live_timing-parser/pkg/helpers"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/tealeg/xlsx/v3"
)

func ReadResultsJson(jsonStr string) (*helpers.ACSMResultsList, error) {
	result := helpers.ACSMResultsList{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	return &result, err
}

func GetResults(jsonStr string) (*helpers.ACSMResults, error) {
	result := helpers.ACSMResults{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	return &result, err
}

func retrieveBestLap(results []helpers.ACSMResult) int {
	bestLap := 999999999
	for _, result := range results {
		if result.BestLap < bestLap {
			bestLap = result.BestLap
		}
	}

	return bestLap
}

func insertExcelCell(row *xlsx.Row, value string) *xlsx.Cell {
	cell := row.AddCell()
	cell.Value = value
	return cell
}

func calculateConsistency(laps []helpers.ACSMLap) map[string]float64 {
	result := map[string]float64{}
	diversLaps := map[string][]int{}
	for _, lap := range laps {
		diversLaps[lap.DriverGUID] = append(diversLaps[lap.DriverGUID], lap.LapTime)
	}

	for driver, laps := range diversLaps {
		data := stats.LoadRawData(laps)
		//https://foro.rinconmatematico.com/index.php?topic=109541.0
		//deviation, err := data.MedianAbsoluteDeviation()
		deviation, err := data.StandardDeviation()
		if err != nil {
			continue
		}
		mean, err := data.Mean()
		if err != nil {
			continue
		}
		result[driver] = deviation / mean
	}
	return result
}

func GenerateResultsExcel(results *helpers.ACSMResults, name string, points []int) (*xlsx.File, error) {

	file := xlsx.NewFile()
	sheet, err := file.AddSheet(name)
	if err != nil {
		return nil, err
	}
	row := sheet.AddRow()
	for _, value := range []string{"Pos", "Points", "Driver", "Laps", "Time/Retired", "Best lap", "Ballast (Kg)", "Gained/Lost", "Consistency", "Penalty"} {
		cell := row.AddCell()
		cell.Value = value
	}

	bestLap := retrieveBestLap(results.Result)
	consistencies := calculateConsistency(results.Laps)

	totalTimeFirst := 0
	totalLapsFirst := 0

	for i, result := range results.Result {
		if i == 0 {
			totalTimeFirst = result.TotalTime
			totalLapsFirst = result.NumLaps
		}

		row = sheet.AddRow()
		insertExcelCell(row, strconv.Itoa(i+1))

		pts := 0
		if i < len(points) {
			pts = points[i]
		}
		if result.BestLap == bestLap {
			pts++
		}
		insertExcelCell(row, strconv.Itoa(pts))

		insertExcelCell(row, result.DriverName)
		insertExcelCell(row, strconv.Itoa(result.NumLaps))

		totalTimeStr := ""
		if result.NumLaps < totalLapsFirst {
			totalTimeStr = fmt.Sprintf("+%d laps", totalLapsFirst-result.NumLaps)
		} else {
			totalTimeInt := totalTimeFirst
			if i != 0 {
				totalTimeInt = result.TotalTime - totalTimeInt
			}
			totalTimeStr = helpers.ConvertTimeToHuman(time.Duration(totalTimeInt * int(time.Millisecond)))
			if i != 0 {
				totalTimeStr = "+" + totalTimeStr
			}
		}
		insertExcelCell(row, totalTimeStr)

		bestLapStr := "-"
		if result.BestLap != 999999999 {
			bestLapStr = helpers.ConvertTimeToHuman(time.Duration(result.BestLap * int(time.Millisecond)))
		}
		insertExcelCell(row, bestLapStr)

		insertExcelCell(row, strconv.Itoa(result.BallastKG))
		insertExcelCell(row, strconv.Itoa(result.GridPosition-(i+1)))

		consistency := "-"
		if c, ok := consistencies[result.DriverGUID]; ok && c > 0 {
			consistency = fmt.Sprintf("%.2f%%", 100-(c*100))
		}
		insertExcelCell(row, consistency)

		penalty := ""
		if result.PenaltyTime > 0 {
			penalty = fmt.Sprintf("+%dsec", result.PenaltyTime/1000000000)
		}
		insertExcelCell(row, penalty)
	}

	return file, nil
}
