package main

import (
	"flag"
	"fmt"
	"github.com/getsentry/sentry-go"
	"go-sentry-cmd/models"
	"os/exec"
	"strings"
	"time"
)

func main() {
	var config = models.Config{}
	parseFlags(&config)
	err := sentry.Init(sentry.ClientOptions{Dsn: config.Dsn, Environment: config.Environment})
	if err != nil {
		panic(err)
	}
	args := strings.Split(config.Command, " ")
	cmd1 := args[0]
	args = append(args[:0], args[0+1:]...)
	cmd := exec.Command(cmd1, args...)
	cmdErr := cmd.Run()
	if cmdErr != nil {
		sentry.CaptureException(cmdErr)
	} else if config.ReportAll {
		sentry.CaptureMessage(fmt.Sprintf("Command \"%v\" completed successfully", config.Command))
	}
	sentry.Flush(time.Second * 5)
}

func parseFlags(config *models.Config) {
	dsn := flag.String("dsn", "null", "Sentry DSN url")
	command := flag.String("command", "null", "Command to run")
	env := flag.String("environment", "undefined", "Sentry environment")
	reportAll := flag.Bool("reportAll", false, "Report success command")
	flag.Parse()

	if *command == "null" {
		panic("Command cant' be nil")
	}
	if *dsn == "null" {
		panic("Dsn can't be nil")
	}
	config.Dsn = *dsn
	config.Command = *command
	config.Environment = *env
	config.ReportAll = *reportAll
}
