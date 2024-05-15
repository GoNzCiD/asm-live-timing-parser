package acsm_parser

import (
	"acsm-live_timing-parser/pkg/helpers"
	"encoding/json"
)

func ReadResultsJson(jsonStr string) (*helpers.ACSMResultsList, error) {
	result := helpers.ACSMResultsList{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	return &result, err
}
