package config

import (
	"time"

	"github.com/hydronica/toml"
)

const (
	DefaultConfigFile = "config.toml"
)

type ServerConfig struct {
	Timeout          time.Duration `toml:"timeout"`
	AllowedClientIPs []string      `toml:"allowed_client_ip_ranges"`
	LogPath          string        `toml:"log_path"`
	HttpLogPath      string        `toml:"http_log_path"`
	Address          string        `toml:"address"`
	Scrape           ScrapeConfig  `toml:"scrape"`
	Race             RaceConfig    `toml:"race"`
}

type ScrapeConfig struct {
	ACSMDomain         string        `toml:"acsm_domain"`
	User               string        `toml:"user"`
	Password           string        `toml:"password"`
	ChromeDriverPath   string        `toml:"chrome-driver_path"`
	SeleniumUrl        string        `toml:"selenium_url"`
	WorkingPath        string        `toml:"working_path"`
	LeaderBoardJsonTtl time.Duration `toml:"leaderboard_json_ttl"`
}

type RaceConfig struct {
	Points               []int  `toml:"points"`
	Ballast              []int  `toml:"ballast"`
	RPLogsPath           string `toml:"rp_logs_path"`
	ClassificationIdsUrl string `toml:"classification_ids_url"`
}

func LoadConfig(cfg *ServerConfig, fileName string) error {
	_, err := toml.DecodeFile(fileName, cfg)
	if err != nil {
		return err
	}

	if len(cfg.AllowedClientIPs) == 0 {
		cfg.AllowedClientIPs = []string{"127.0.0.0/8", "::1/128"}
	}

	return nil
}
