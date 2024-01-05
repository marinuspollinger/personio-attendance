package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	config "personio-cli/config"
	"time"
)

type personioServerResponse struct {
	StartTime   time.Time `json:"start_time"`
	CurrentTime time.Time `json:"current_time"`
	BreakStart  time.Time `json:"break_start"`
	BreakEnd    time.Time `json:"break_end"`
	Success     bool      `json:"success"`
	Error       string    `json:"error_message"`
}

type setCustomTimeRequest struct {
	BreakStartTime time.Time `json:"break_start_time"`
	BreakStopTime  time.Time `json:"break_end_time"`
	StartTime      time.Time `json:"start_time"`
	StopTime       time.Time `json:"end_time"`
	Override       bool      `json:"override"`
}

func GetStatus(cfg config.EnvConfig) (int, error) {
	err := getRequest("/api/getServerData", cfg)
	if err != nil {
		return 1, err
	}
	return 0, nil
}

func SetBreakStartTime(cfg config.EnvConfig, start string, override bool) (int, error) {
	startParsed, err := parseTimeStamp(start)
	if err != nil {
		return 1, err
	}

	body := setCustomTimeRequest{
		BreakStartTime: startParsed,
		Override:       override,
	}

	err = postRequest("/api/setBreakStartTime", body, cfg)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

func SetBreakStopTime(cfg config.EnvConfig, stop string, override bool) (int, error) {
	stopParsed, err := parseTimeStamp(stop)
	if err != nil {
		return 1, err
	}

	body := setCustomTimeRequest{
		BreakStopTime: stopParsed,
		Override:      override,
	}

	err = postRequest("/api/setBreakStopTime", body, cfg)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

func SetStartTime(cfg config.EnvConfig, start string) (int, error) {
	startParsed, err := parseTimeStamp(start)
	if err != nil {
		return 1, err
	}

	body := setCustomTimeRequest{
		StartTime: startParsed,
	}

	err = postRequest("/api/setStartTime", body, cfg)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

func SetStopTime(cfg config.EnvConfig, stop string) (int, error) {
	stopParsed, err := parseTimeStamp(stop)
	if err != nil {
		return 1, err
	}

	body := setCustomTimeRequest{
		StopTime: stopParsed,
	}

	err = postRequest("/api/setStopTime", body, cfg)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

func WriteTimes(cfg config.EnvConfig, yFlag bool) (int, error) {
	var selection string

	_, err := GetStatus(cfg)
	if err != nil {
		return 1, err
	}

	// when -y is not set, print out current times
	if !yFlag {
		fmt.Println("are you sure, you want to proceed? (y/n)")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		selection = input.Text()

	}

	if selection == "y" || selection == "Y" || yFlag {
		fmt.Printf("\n\nSending Request to Personio...\n\n")
		err := getRequest("/api/sendToPersonio", cfg)
		if err != nil {
			return 1, err
		}
		return 0, nil
	}

	fmt.Println("Aborting...")
	return 1, nil
}

func getRequest(endpoint string, cfg config.EnvConfig) error {
	res, err := http.Get("http://" + cfg.HttpAddress + endpoint)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var data personioServerResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	if endpoint == "/api/getServerData" {
		printToShell(data)
	}

	if !data.Success {
		return fmt.Errorf(data.Error)
	}

	return nil

}

func postRequest(endpoint string, body setCustomTimeRequest, cfg config.EnvConfig) error {
	jsonbody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	res, err := http.Post("http://"+cfg.HttpAddress+endpoint, "application/json", bytes.NewBuffer(jsonbody))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resbody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var data personioServerResponse
	if err := json.Unmarshal(resbody, &data); err != nil {
		return err
	}

	if !data.Success {
		return fmt.Errorf(data.Error)
	}

	return nil
}

func printToShell(data personioServerResponse) {
	fmt.Printf("Personio CLI Tool by Marinus\n")
	fmt.Printf("Start-Time (logon time): \t\t\t\t%s\n", data.StartTime.Format(config.ParseTimeStampString))
	fmt.Printf("Current-Time (End time when sending to personio): \t%s\n", data.CurrentTime.Format(config.ParseTimeStampString))
	fmt.Printf("Break-Start-Time: \t\t\t\t\t%s\n", data.BreakStart.Format(config.ParseTimeStampString))
	fmt.Printf("Break-End-Time: \t\t\t\t\t%s\n", data.BreakEnd.Format(config.ParseTimeStampString))
	fmt.Printf("Last Server Response Success: \t\t\t\t%v\n", data.Success)
	fmt.Printf("Last Server Response Error: \t\t\t\t%s\n", data.Error)
}

func parseTimeStamp(ts string) (time.Time, error) {
	loc, err := time.LoadLocation("CET")
	if err != nil {
		return time.Time{}, err
	}

	parsedTs, err := time.ParseInLocation(config.ParseTimeStampString, ts, loc)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTs, nil
}
