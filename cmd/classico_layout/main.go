package main

import (
	syslog "log"
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/bhmj/classico_layout/internal/app/app"
	"github.com/bhmj/classico_layout/internal/pkg/config"
	"github.com/bhmj/classico_layout/internal/pkg/log"
)

func main() {
	cfg := config.New()
	parser := flags.NewParser(cfg, flags.Default)
	parser.LongDescription = "Classico tile layout generator\nversion 0.3.0"
	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	if cfg.ConfigFile != "" {
		err = cfg.ReadFromFile(cfg.ConfigFile)
		if err != nil {
			syslog.Fatal(err.Error())
		}
	}

	logger, err := log.New(cfg.LogLevel)
	if err != nil {
		syslog.Fatal(err.Error())
	}
	defer func() { _ = logger.L().Sync() }()

	app.New(cfg, logger).Run()
}
