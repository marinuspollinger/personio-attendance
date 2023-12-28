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

	status := flag.Bool("status", false, "Get Current Times")
	sendToPersonio := flag.Bool("send", false, "Send Current Times to Personio API")
	yFlag := flag.Bool("yes", false, "Immediatly write to personio, without checking times first")

	startBreak := flag.Bool("start-break", false, "Start the Break now!")
	endBreak := flag.Bool("end-break", false, "End the Break now!")

	customBreakStart := flag.String("break-start-time", "", "Set Break Start Time")
	customBreakStop := flag.String("break-end-time", "", "Set Break End Time")

	flag.Parse()

	// both flags MUST be set at the same time if they are used
	if *customBreakStart != "" || *customBreakStop != "" {
		if (*customBreakStart == "") != (*customBreakStop == "") {
			fmt.Println("Both start-time and end-time must be set together")
			os.Exit(1)
		}
	}

	// check if more then one flag are set, only check one customBreak flag because of previous check above
	flagsSet := 0
	if *status {
		flagsSet++
	}
	if *startBreak {
		flagsSet++
	}
	if *endBreak {
		flagsSet++
	}
	if *sendToPersonio {
		flagsSet++
	}
	if *customBreakStart != "" {
		flagsSet++
	}

	// if more then one flag are set, fail
	if flagsSet > 1 {
		fmt.Println("Only one flag can be set at a time")
		os.Exit(1)
	}

	// print help message if no flag is set
	if flagsSet == 0 {
		flag.PrintDefaults()
	}

	if *status {
		exitCode, err := cli.GetStatus(cfg)
		if err != nil {
			fmt.Printf("Error getting status: %s", err)
		}
		os.Exit(exitCode)
	}

	if *startBreak {
		exitCode, err := cli.StartBreak(cfg)
		if err != nil {
			fmt.Printf("Error starting Break: %s", err)
		}
		os.Exit(exitCode)

	}

	if *endBreak {
		exitCode, err := cli.EndBreak(cfg)
		if err != nil {
			fmt.Printf("Error ending Break: %s", err)
		}
		os.Exit(exitCode)

	}

	if *customBreakStart != "" {
		exitCode, err := cli.CustomBreakTimes(cfg, *customBreakStart, *customBreakStop)
		if err != nil {
			fmt.Printf("Error setting Custom Break Times: %s", err)
		}
		os.Exit(exitCode)
	}

	if *sendToPersonio {
		exitCode, err := cli.WriteTimes(cfg, *yFlag)
		if err != nil {
			fmt.Printf("Error sending times to Personio: %s", err)
		}
		os.Exit(exitCode)
	}

}
