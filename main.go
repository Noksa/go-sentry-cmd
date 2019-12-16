package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/Noksa/go-sentry-cmd/models"
	"github.com/getsentry/sentry-go"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	var config = models.Config{}
	parseFlags(&config)
	err := sentry.Init(sentry.ClientOptions{Dsn: config.Dsn, Environment: config.Environment, AttachStacktrace: true})
	if err != nil {
		panic(err)
	}
	args := strings.Split(config.Command, " ")
	cmd1 := args[0]
	args = append(args[:0], args[0+1:]...)
	cmd := exec.Command(cmd1, args...)
	var outBuffer = bytes.Buffer{}
	var errBuffer = bytes.Buffer{}
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer
	cmdErr := cmd.Run()
	SetAdditionalInfo()
	if cmdErr != nil {
		var errorMsg = fmt.Sprintf("Command \"%v\" completed with errors!\nResult: %v\nAdditional data: %v", config.Command, cmdErr.Error(), errBuffer.String())
		var newErr = errors.New(errorMsg)
		sentry.CaptureException(newErr)
	} else if config.ReportAll {
		var res = outBuffer.String()
		sentry.CaptureMessage(fmt.Sprintf("Command \"%v\" completed successfully!\nResult: %v", config.Command, res))
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
		panic("Command cant' be nil or empty")
	}

	if *env == "undefined" {
		osEnv := os.Getenv("SENTRY_ENVIRONMENT")
		if osEnv != "" {
			*env = osEnv
		}
	}

	if *dsn == "null" {
		*dsn = os.Getenv("SENTRY_DSN")
		if *dsn == "" {
			panic("Dsn can't be nil or empty")
		}
	}
	config.Dsn = *dsn
	config.Command = *command
	config.Environment = *env
	config.ReportAll = *reportAll
}

func SetAdditionalInfo() {
	hostName := os.Getenv("SENTRY_HOSTNAME")
	if hostName != "" {
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("server_name", hostName)
		})
	}
}
