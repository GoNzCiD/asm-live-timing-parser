package helpers

import (
	"time"
)

type ACSMResultsDriver struct {
	GUID      string   `json:"Guid"`
	GuidsList []string `json:"GuidsList"`
	Name      string   `json:"Name"`
	Nation    string   `json:"Nation"`
	Team      string   `json:"Team"`
	ClassID   string   `json:"ClassID"`
}

type ACSMResults struct {
	Version int `json:"Version"`
	Cars    []struct {
		BallastKG  int               `json:"BallastKG"`
		CarID      int               `json:"CarId"`
		Driver     ACSMResultsDriver `json:"Driver"`
		Model      string            `json:"Model"`
		Restrictor int               `json:"Restrictor"`
		Skin       string            `json:"Skin"`
		ClassID    string            `json:"ClassID"`
		MinPing    int               `json:"MinPing"`
		MaxPing    int               `json:"MaxPing"`
	} `json:"Cars"`
	Events []struct {
		CarID           int               `json:"CarId"`
		Driver          ACSMResultsDriver `json:"Driver"`
		ImpactSpeed     float64           `json:"ImpactSpeed"`
		OtherCarID      int               `json:"OtherCarId"`
		OtherDriver     ACSMResultsDriver `json:"OtherDriver"`
		RelPosition     Position          `json:"RelPosition"`
		Type            string            `json:"Type"`
		WorldPosition   Position          `json:"WorldPosition"`
		Timestamp       int               `json:"Timestamp"`
		AfterSessionEnd bool              `json:"AfterSessionEnd"`
	} `json:"Events"`
	Laps []struct {
		BallastKG               int         `json:"BallastKG"`
		CarID                   int         `json:"CarId"`
		CarModel                string      `json:"CarModel"`
		Cuts                    int         `json:"Cuts"`
		DriverGUID              string      `json:"DriverGuid"`
		DriverName              string      `json:"DriverName"`
		LapTime                 int         `json:"LapTime"`
		Restrictor              int         `json:"Restrictor"`
		Sectors                 []int       `json:"Sectors"`
		Timestamp               int         `json:"Timestamp"`
		Tyre                    string      `json:"Tyre"`
		ClassID                 string      `json:"ClassID"`
		ContributedToFastestLap bool        `json:"ContributedToFastestLap"`
		SpeedTrapHits           interface{} `json:"SpeedTrapHits"`
		Conditions              struct {
			Ambient       float64 `json:"Ambient"`
			Road          float64 `json:"Road"`
			Grip          int     `json:"Grip"`
			WindSpeed     float64 `json:"WindSpeed"`
			WindDirection float64 `json:"WindDirection"`
			RainIntensity int     `json:"RainIntensity"`
			RainWetness   int     `json:"RainWetness"`
			RainWater     int     `json:"RainWater"`
		} `json:"Conditions"`
	} `json:"Laps"`
	Result []struct {
		BallastKG    int    `json:"BallastKG"`
		BestLap      int    `json:"BestLap"`
		CarID        int    `json:"CarId"`
		CarModel     string `json:"CarModel"`
		DriverGUID   string `json:"DriverGuid"`
		DriverName   string `json:"DriverName"`
		Restrictor   int    `json:"Restrictor"`
		TotalTime    int    `json:"TotalTime"`
		NumLaps      int    `json:"NumLaps"`
		HasPenalty   bool   `json:"HasPenalty"`
		PenaltyTime  int    `json:"PenaltyTime"`
		LapPenalty   int    `json:"LapPenalty"`
		Disqualified bool   `json:"Disqualified"`
		ClassID      string `json:"ClassID"`
		GridPosition int    `json:"GridPosition"`
	} `json:"Result"`
	Penalties     interface{} `json:"Penalties"`
	TrackConfig   string      `json:"TrackConfig"`
	TrackName     string      `json:"TrackName"`
	Type          string      `json:"Type"`
	Date          time.Time   `json:"Date"`
	SessionFile   string      `json:"SessionFile"`
	SessionConfig struct {
		SessionType                     int    `json:"session_type"`
		Name                            string `json:"name"`
		Time                            int    `json:"time"`
		Laps                            int    `json:"laps"`
		IsOpen                          int    `json:"is_open"`
		WaitTime                        int    `json:"wait_time"`
		VisibilityMode                  int    `json:"visibility_mode"`
		QualifyingType                  int    `json:"qualifying_type"`
		QualifyingNumberOfLapsToAverage int    `json:"qualifying_number_of_laps_to_average"`
		CountOutLap                     bool   `json:"count_out_lap"`
		DisablePushToPass               bool   `json:"disable_push_to_pass"`
	} `json:"SessionConfig"`
	ChampionshipID string `json:"ChampionshipID"`
	RaceWeekendID  string `json:"RaceWeekendID"`
}

type ACSMResultsList struct {
	NumPages    int    `json:"num_pages"`
	CurrentPage int    `json:"current_page"`
	SortType    string `json:"sort_type"`
	Results     []struct {
		Track          string    `json:"track"`
		TrackLayout    string    `json:"track_layout,omitempty"`
		SessionType    string    `json:"session_type"`
		Date           time.Time `json:"date"`
		ResultsJSONURL string    `json:"results_json_url"`
		ResultsPageURL string    `json:"results_page_url"`
	} `json:"results"`
}
