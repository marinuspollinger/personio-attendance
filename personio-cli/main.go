package main

import (
	"flag"
	"fmt"
	"os"
	"personio-cli/cli"
	config "personio-cli/config"

	"github.com/codingconcepts/env"
	ulog "github.com/dunv/ulog/v2"
)

func main() {
	cfg := config.EnvConfig{}
	ulog.FatalIfError(env.Set(&cfg))

	sendToPersonio := flag.Bool("send", false, "Send Current Times to Personio API")
	yFlag := flag.Bool("yes", false, "Immediatly write to personio, without checking times first")

	startBreakNow := flag.Bool("start-break", false, "Start the Break now!")
	endBreakNow := flag.Bool("end-break", false, "End the Break now!")

	customStartTime := flag.String("start-time", "", "Set Start Time")
	customStopTime := flag.String("stop-time", "", "Set Stop Time")
	customBreakStart := flag.String("break-start-time", "", "Set Break Start Time")
	customBreakStop := flag.String("break-stop-time", "", "Set Break Stop Time")

	flag.Parse()

	// check if more then one flag are set, only check one customBreak flag because of previous check above
	uniqueFlagsSet := 0
	if *startBreakNow {
		uniqueFlagsSet++
	}
	if *endBreakNow {
		uniqueFlagsSet++
	}
	if *sendToPersonio {
		uniqueFlagsSet++
	}

	// if more then one flag are set, fail
	if uniqueFlagsSet > 1 {
		fmt.Println("Only one flag can be set at a time")
		os.Exit(1)
	}

	if *startBreakNow {
		exitCode, err := cli.StartBreak(cfg)
		if err != nil {
			fmt.Printf("Error starting Break: %s", err)
		}
		exit(cfg, exitCode)

	}

	if *endBreakNow {
		exitCode, err := cli.EndBreak(cfg)
		if err != nil {
			fmt.Printf("Error ending Break: %s", err)
		}
		exit(cfg, exitCode)

	}

	if *customStartTime != "" {
		exitCode, err := cli.CustomStartTime(cfg, *customStartTime)
		if err != nil {
			fmt.Printf("Error setting Custom Start Time: %s", err)
		}
		if exitCode != 0 {
			exit(cfg, exitCode)
		}
	}

	if *customBreakStart != "" {
		exitCode, err := cli.CustomBreakStartTime(cfg, *customBreakStart)
		if err != nil {
			fmt.Printf("Error setting Custom Break Start Time: %s", err)
		}
		if exitCode != 0 {
			exit(cfg, exitCode)
		}
	}

	if *customBreakStop != "" {
		exitCode, err := cli.CustomBreakStopTime(cfg, *customBreakStop)
		if err != nil {
			fmt.Printf("Error setting Custom Break Stop Time: %s", err)
		}
		if exitCode != 0 {
			exit(cfg, exitCode)
		}
	}

	if *customStopTime != "" {
		exitCode, err := cli.CustomStopTime(cfg, *customStopTime)
		if err != nil {
			fmt.Printf("Error setting Custom Stop Time: %s", err)
		}
		if exitCode != 0 {
			exit(cfg, exitCode)
		}
	}

	if *sendToPersonio {
		exitCode, err := cli.WriteTimes(cfg, *yFlag)
		if err != nil {
			fmt.Printf("Error sending times to Personio: %s", err)
		}
		if exitCode != 0 {
			exit(cfg, exitCode)
		}
	}

	exit(cfg, 0)
}

func exit(cfg config.EnvConfig, exitCode int) {
	_, err := cli.GetStatus(cfg)
	if err != nil {
		fmt.Printf("Error getting status: %s", err)
		os.Exit(127)
	}

	os.Exit(exitCode)
}
