package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"reverseproxy/internal/config"
	masks "reverseproxy/internal/masker"
	"reverseproxy/proxy"

	"github.com/rs/zerolog"
)

var (
	tomlPathFlag = flag.String("config", "./internal/config/example_config.toml",
		"Specify the path of config.toml file, e.g.: -config /folder/config.toml")
)

func main() {
	flag.Parse()
	// init logger
	zerolog.DurationFieldInteger = true
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	// parse config
	cfg, err := config.LoadConfig(*tomlPathFlag)
	if err != nil {
		log.Panic().Err(err).Msg("failed to load config")
	}
	// Create Blockers from config
	blockers := addBlockersFromConfig(cfg)
	// Create Masker
	masker := []proxy.Masker{
		masks.NewEmailMasker(),
		masks.NewCreditCardMasker()}
	// Create Proxy
	rp, err := proxy.New(
		cfg.TargetURL,
		cfg.ReverseProxyPort,
		masker,
		blockers,
		log)
	if err != nil {
		log.Panic().Err(err).Msg("failed to create reverse proxy")
	}
	// bind signals to quit channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Start it
	cancel, err := rp.Start()
	if err != nil {
		log.Panic().Err(err).Msg("failed to start reverse proxy")
	}
	log.Info().Msg("reverse proxy started")
	// Wait for quit signal
	<-quit
	// Gracefully shutdown reverse proxy
	cancel()
	time.Sleep(time.Second * 1)
	log.Info().Msg("reverse proxy stopped")
}

func addBlockersFromConfig(cfg *config.Config) []proxy.Blocker {
	var blockers []proxy.Blocker
	if cfg.HeaderBlocker != nil {
		blockers = append(blockers, cfg.HeaderBlocker)
	}
	if cfg.ParamBlocker != nil {
		blockers = append(blockers, cfg.ParamBlocker)
	}
	if cfg.PathBlocker != nil {
		blockers = append(blockers, cfg.PathBlocker)
	}
	if cfg.MethodBlocker != nil {
		blockers = append(blockers, cfg.MethodBlocker)
	}
	return blockers
}
