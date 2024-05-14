package main

import (
	downloader "acsm-live_timing-parser/pkg/downloader"
	"acsm-live_timing-parser/pkg/helpers"
	"acsm-live_timing-parser/pkg/leaderboard_parser"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

var Opts struct {
	chromeDriverPath       string
	seleniumPath           string
	debug                  bool
	user                   string
	pass                   string
	url                    string
	serverNo               int
	live                   bool
	outputPath             string
	oldJsonTempPathSeconds int
	skinPreviewPattern     string
}

func parseArgs() {
	flag.StringVar(&Opts.chromeDriverPath, "c", "", "Chrome driver path")
	flag.StringVar(&Opts.seleniumPath, "w", "", "Selenium Chrome URL")
	flag.StringVar(&Opts.url, "b", "", "Base home URL for Assetto Corsa Server Manager")
	flag.IntVar(&Opts.serverNo, "i", 0, "Instance server (for multiserver ACSM)")
	flag.StringVar(&Opts.user, "u", "", "ACSM credentials user")
	flag.StringVar(&Opts.pass, "p", "", "ACSM credentials password")
	flag.BoolVar(&Opts.debug, "d", false, "Open browser window")
	flag.BoolVar(&Opts.live, "l", false, "Do not use json temporal file")
	flag.StringVar(&Opts.outputPath, "o", "livetime.json", "Output json live time path")
	flag.IntVar(&Opts.oldJsonTempPathSeconds, "t", helpers.DefaultOldJsonTempPathSeconds, "Time in seconds between last API call")
	flag.StringVar(&Opts.skinPreviewPattern, "s", "", "URL path pattern for Skin image preview")

	flag.Parse()

	Opts.chromeDriverPath, _ = helpers.GetFullPath(Opts.chromeDriverPath)
	Opts.outputPath, _ = helpers.GetFullPath(Opts.outputPath)
}

func callApi(tmpPath string) (string, error) {
	downloader := downloader.NewDownloader(Opts.chromeDriverPath, Opts.seleniumPath, Opts.user, Opts.pass, Opts.url, Opts.serverNo, Opts.debug)
	jsonStr, err := downloader.Download()
	if err != nil {
		return "", err
	}
	err = helpers.SaveToFile(jsonStr, tmpPath)
	if err != nil {
		fmt.Printf("cannot save to the temporal json file: %v\n", err)
	}
	return jsonStr, nil
}

func main() {
	parseArgs()

	var jsonStr string

	jsonPath, err := helpers.GetFullPath(helpers.TmpLeaderBoardFileName)
	if err != nil {
		fmt.Printf("cannot retrieve full path for temporal json file: %v\n", err)
		os.Exit(1)
	}

	info, err := os.Stat(jsonPath)
	if err != nil {
		fmt.Printf("cannot retrieve the temporal json file stats: %v\n", err)
	}
	if err != nil || time.Since(info.ModTime()) > time.Duration(Opts.oldJsonTempPathSeconds)*time.Second {
		jsonStr, err = callApi(jsonPath)
		if err != nil {
			fmt.Printf("cannot download data: %v\n", err)
			os.Exit(1)
		}
	}

	if jsonStr == "" {
		jsonStr, err = helpers.ReadFromFile(jsonPath)
		if err != nil {
			jsonStr, err = callApi(jsonPath)
			if err != nil {
				fmt.Printf("cannot download data: %v\n", err)
				os.Exit(1)
			}
		}
	}

	liveTiming, err := leaderboard_parser.ReadJson(jsonStr)
	if err != nil {
		fmt.Printf("cannot parse API json value: %v\n", err)
		os.Exit(1)
	}

	result := leaderboard_parser.ExtractHotlaps(liveTiming.ConnectedDrivers)
	result = append(result, leaderboard_parser.ExtractHotlaps(liveTiming.DisconnectedDrivers)...)
	result = leaderboard_parser.SortAndCalculateData(result, &Opts.skinPreviewPattern)

	resultBytes, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		fmt.Printf("cannot format output data: %v\n", err)
		os.Exit(1)
	}

	err = helpers.SaveToFile(string(resultBytes), Opts.outputPath)
	if err != nil {
		fmt.Printf("cannot write to the output file (%s): %v\n", Opts.outputPath, err)
		os.Exit(1)
	}
}
