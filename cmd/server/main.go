package main

import (
	"acsm-live_timing-parser/internal/server"
	"acsm-live_timing-parser/pkg/config"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type ProgramOpts struct {
	ConfigFile string
	UrlPrefix  string
	Verbose    bool
}

var Opts ProgramOpts
var Cfg config.ServerConfig

func ParseArgs() {
	flag.StringVar(&Opts.ConfigFile, "c", config.DefaultConfigFile, "Configuration file")
	flag.StringVar(&Opts.UrlPrefix, "p", "/", "URL path prefix")
	flag.BoolVar(&Opts.Verbose, "v", false, "Set logging level to DEBUG, and ignore setting from config")
	flag.Parse()

	Opts.UrlPrefix = strings.TrimRight(Opts.UrlPrefix, "/")
}

func main() {
	ParseArgs()

	level := slog.LevelInfo
	if Opts.Verbose {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})))

	err := config.LoadConfig(&Cfg, Opts.ConfigFile)
	if err != nil {
		slog.Error(fmt.Sprintf("could not load config file: %v", err.Error()))
		os.Exit(1)
	}

	if Cfg.LogPath != "" {
		f, err := os.OpenFile(Cfg.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			slog.Error(fmt.Sprintf("could not open log file: %v", err.Error()))
		} else {
			slog.SetDefault(slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{Level: level})))
		}
		defer f.Close()
	}

	var httpLog *os.File
	if Cfg.HttpLogPath != "" {
		f, err := os.OpenFile(Cfg.HttpLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			slog.Error(fmt.Sprintf("could not open http log file: %v", err.Error()))
		} else {
			httpLog = f
		}
		defer f.Close()
	}

	s := server.NewServer(httpLog, Cfg.Address, Opts.UrlPrefix, &Cfg.ScrapeConfig)
	slog.Debug("Initializing server")
	err = s.InitializeAndStart()
	if err != nil {
		slog.Error(fmt.Sprintf("cannot start the server: %v", err))
		os.Exit(1)
	}
}
