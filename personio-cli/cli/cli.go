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

type setBreakCustomTimeRequest struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func GetStatus(cfg config.EnvConfig) (int, error) {

	err := getRequest("/api/getServerData", cfg)
	if err != nil {
		return 1, err
	}
	return 0, nil
}

func EndBreak(cfg config.EnvConfig) (int, error) {
	err := getRequest("/api/endBreakNow", cfg)
	if err != nil {
		return 1, err
	}
	return 0, nil
}

func StartBreak(cfg config.EnvConfig) (int, error) {
	err := getRequest("/api/startBreakNow", cfg)
	if err != nil {
		return 1, err
	}
	return 0, nil
}

func CustomBreakTimes(cfg config.EnvConfig, start string, stop string) (int, error) {
	parseTimeStamp := "2006-01-02 15:04"
	loc, err := time.LoadLocation("CET")
	if err != nil {
		return 1, err
	}

	startParsed, err := time.ParseInLocation(parseTimeStamp, start, loc)
	if err != nil {
		return 1, err
	}

	endParsed, err := time.ParseInLocation(parseTimeStamp, stop, loc)
	if err != nil {
		return 1, err
	}

	body := setBreakCustomTimeRequest{
		StartTime: startParsed,
		EndTime:   endParsed,
	}

	err = postRequest("/api/setBreakCustomTime", body, cfg)
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

	printToShell(data)
	return nil

}

func postRequest(endpoint string, body setBreakCustomTimeRequest, cfg config.EnvConfig) error {
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

	printToShell(data)
	return nil
}

func printToShell(data personioServerResponse) {
	fmt.Printf("Personio CLI Tool by Marinus\n")
	fmt.Printf("Start-Time (logon time): \t\t%s\n", data.StartTime.Format(time.RFC822))
	fmt.Printf("Current-Time: \t\t\t\t%s\n", data.CurrentTime.Format(time.RFC822))
	fmt.Printf("Break-Start-Time: \t\t\t%s\n", data.BreakStart.Format(time.RFC822))
	fmt.Printf("Break-End-Time: \t\t\t%s\n", data.BreakEnd.Format(time.RFC822))
	fmt.Printf("Local-Server Response Success: \t\t%v\n", data.Success)
	fmt.Printf("Local-Server Response Error: \t\t%s\n", data.Error)
}
