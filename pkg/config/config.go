package config

import (
	"time"

	"github.com/hydronica/toml"
)

type ServerConfig struct {
	Timeout          time.Duration `toml:"timeout"`
	AllowedClientIPs []string      `toml:"allowed_client_ip_ranges"`
	LogPath          string        `toml:"log_path"`
	HttpLogPath      string        `toml:"http_log_path"`
	Address          string        `toml:"address"`
	ScrapeConfig     ScrapeConfig  `toml:"scrape"`
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

const (
	DefaultConfigFile = "config.toml"
)

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
