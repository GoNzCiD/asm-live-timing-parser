package helpers

import (
	"time"
)

type CarInfo struct {
	CarID          int    `json:"CarID"`
	DriverName     string `json:"DriverName"`
	DriverGUID     string `json:"DriverGUID"`
	CarModel       string `json:"CarModel"`
	CarSkin        string `json:"CarSkin"`
	Tyres          string `json:"Tyres"`
	RaceNumber     int    `json:"RaceNumber"`
	DriverInitials string `json:"DriverInitials"`
	CarName        string `json:"CarName"`
	ClassID        string `json:"ClassID"`
	EventType      int    `json:"EventType"`
}

type Position struct {
	X float64 `json:"X"`
	Y float64 `json:"Y"`
	Z float64 `json:"Z"`
}

type DriverSwap struct {
	IsInDriverSwap                   bool      `json:"IsInDriverSwap"`
	DriverSwapExpectedCompletionTime time.Time `json:"DriverSwapExpectedCompletionTime"`
	DriverSwapHadPenalty             bool      `json:"DriverSwapHadPenalty"`
	DriverSwapActualCompletionTime   time.Time `json:"DriverSwapActualCompletionTime"`
}

type DriftInfo struct {
	CurrentLapScore int `json:"CurrentLapScore"`
	BestLapScore    int `json:"BestLapScore"`
}

type Split struct {
	SplitIndex    int           `json:"SplitIndex"`
	SplitTime     time.Duration `json:"SplitTime"`
	Cuts          int           `json:"Cuts"`
	IsDriversBest bool          `json:"IsDriversBest"`
	IsBest        bool          `json:"IsBest"`
}

//func (s *Split) MarshalJSON() ([]byte, error) {
//	return json.Marshal(&struct {
//		SplitTime     string `json:"SplitTime"`
//		Cuts          int    `json:"Cuts"`
//		IsDriversBest bool   `json:"IsDriversBest"`
//		IsBest        bool   `json:"IsBest"`
//	}{
//		SplitTime:     fmt.Sprintf("%.3f", time.Duration(s.SplitTime).Seconds()),
//		Cuts:          s.Cuts,
//		IsDriversBest: s.IsDriversBest,
//		IsBest:        s.IsBest,
//	})
//}

type Splits struct {
	S1 Split `json:"0"`
	S2 Split `json:"1"`
	S3 Split `json:"2"`
}

type MiniSector struct {
	Time      int64 `json:"Time"`
	TimeInLap int   `json:"TimeInLap"`
	Lap       int   `json:"Lap"`
}

type Car struct {
	TopSpeedThisLap      float64      `json:"TopSpeedThisLap"`
	TopSpeedBestLap      float64      `json:"TopSpeedBestLap"`
	TyreBestLap          string       `json:"TyreBestLap"`
	BestLap              int64        `json:"BestLap"`
	NumLaps              int          `json:"NumLaps"`
	LastLap              int64        `json:"LastLap"`
	LastLapCompletedTime time.Time    `json:"LastLapCompletedTime"`
	TotalLapTime         int64        `json:"TotalLapTime"`
	CarName              string       `json:"CarName"`
	RaceNumber           int          `json:"RaceNumber"`
	CurrentLapSplits     Splits       `json:"CurrentLapSplits"`
	BestSplits           Splits       `json:"BestSplits"`
	BestLapSplits        Splits       `json:"BestLapSplits"`
	BestLapMiniSectors   []MiniSector `json:"BestLapMiniSectors"`
}

type Driver struct {
	CarInfo             CarInfo        `json:"CarInfo"`
	TotalNumLaps        int            `json:"TotalNumLaps"`
	ConnectedTime       time.Time      `json:"ConnectedTime"`
	LoadedTime          time.Time      `json:"LoadedTime"`
	Position            int            `json:"Position"`
	Split               string         `json:"Split"`
	DeltaToBest         int            `json:"DeltaToBest"`
	DeltaToSelf         int            `json:"DeltaToSelf"`
	LastSeen            time.Time      `json:"LastSeen"`
	LastPos             Position       `json:"LastPos"`
	IsInPits            bool           `json:"IsInPits"`
	LastPitStop         int            `json:"LastPitStop"`
	DRSActive           bool           `json:"DRSActive"`
	NormalisedSplinePos float64        `json:"NormalisedSplinePos"`
	SteerAngle          int            `json:"SteerAngle"`
	StatusBytes         int            `json:"StatusBytes"`
	BlueFlag            bool           `json:"BlueFlag"`
	HasCompletedSession bool           `json:"HasCompletedSession"`
	Ping                int            `json:"Ping"`
	DriverSwap          DriverSwap     `json:"DriverSwap"`
	DriftInfo           DriftInfo      `json:"DriftInfo"`
	Cars                map[string]Car `json:"Cars"`
}

type LeaderBoard struct {
	Version             int      `json:"Version"`
	SessionIndex        int      `json:"SessionIndex"`
	CurrentSessionIndex int      `json:"CurrentSessionIndex"`
	SessionCount        int      `json:"SessionCount"`
	ServerName          string   `json:"ServerName"`
	Track               string   `json:"Track"`
	TrackConfig         string   `json:"TrackConfig"`
	Name                string   `json:"Name"`
	Type                int      `json:"Type"`
	Time                int      `json:"Time"`
	Laps                int      `json:"Laps"`
	WaitTime            int      `json:"WaitTime"`
	AmbientTemp         int      `json:"AmbientTemp"`
	RoadTemp            int      `json:"RoadTemp"`
	WeatherGraphics     string   `json:"WeatherGraphics"`
	ElapsedMilliseconds int      `json:"ElapsedMilliseconds"`
	VisibilityMode      int      `json:"VisibilityMode"`
	EventType           int      `json:"EventType"`
	ConnectedDrivers    []Driver `json:"ConnectedDrivers"`
	DisconnectedDrivers []Driver `json:"DisconnectedDrivers"`
}
