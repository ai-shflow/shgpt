package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"

	"github.com/cligpt/shgpt/hellogpt"
	"github.com/cligpt/shgpt/hellogpt/config"
)

const (
	gptName = "hellogpt"
)

var (
	app      = kingpin.New(gptName, "hello gpt").Version(config.Version + "-build-" + config.Build)
	logLevel = app.Flag("log-level", "Log level (DEBUG|INFO|WARN|ERROR)").Default("WARN").String()
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	logger, err := initLogger(*logLevel)
	if err != nil {
		return errors.Wrap(err, "failed to init logger")
	}

	gpt, err := initGpt(logger)
	if err != nil {
		return errors.Wrap(err, "failed to init gpt")
	}

	if err := runGpt(logger, gpt); err != nil {
		return errors.Wrap(err, "failed to run gpt")
	}

	return nil
}

func initLogger(level string) (hclog.Logger, error) {
	return hclog.New(&hclog.LoggerOptions{
		Name:  gptName,
		Level: hclog.LevelFromString(level),
	}), nil
}

func initGpt(logger hclog.Logger) (*hellogpt.Context, error) {
	return hellogpt.NewContext(logger.IsDebug())
}

func runGpt(logger hclog.Logger, gpt *hellogpt.Context) error {
	if err := gpt.Init(); err != nil {
		return errors.Wrap(err, "failed to init")
	}

	defer func(gpt *hellogpt.Context) {
		_ = gpt.Deinit()
	}(gpt)

	args := map[string]string{
		hellogpt.Options[0]: "value",
	}

	ret, err := gpt.Run(args)
	if err != nil {
		return errors.Wrap(err, "failed to run")
	}

	logger.Info(ret)

	return nil
}
