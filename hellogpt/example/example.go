package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cligpt/shgpt/hellogpt"
	"github.com/cligpt/shgpt/hellogpt/config"
)

const (
	gptName = "hellogpt"
)

var (
	logLevel string
)

var rootCmd = &cobra.Command{
	Use:     gptName,
	Version: config.Version + "-build-" + config.Build,
	Short:   "hello gpt",
	Long:    "hello gpt",
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(loadConfig(context.Background()))
	},
}

// nolint: gochecknoinits
func init() {
	cobra.OnInitialize()

	rootCmd.Flags().StringVarP(&logLevel, "log-level", "l", "WRAN", "log level (DEBUG|INFO|WARN|ERROR)")
}

func execute() error {
	return rootCmd.Execute()
}

func main() {
	if err := execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func loadConfig(_ context.Context) error {
	logger, err := initLogger(logLevel)
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
