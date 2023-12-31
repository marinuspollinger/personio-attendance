package pserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"personio-server/config"
	"time"

	ulog "github.com/dunv/ulog/v2"
)

type ServerData struct {
	StartTime      time.Time `json:"start_time"`
	StopTime       time.Time `json:"current_time"`
	BreakStart     time.Time `json:"break_start"`
	BreakEnd       time.Time `json:"break_end"`
	Success        bool      `json:"success"`
	Error          string    `json:"error_message"`
	CustomTimeLock bool      `json:"custom_time_lock"`
}

type PersonioRequest struct {
	Attendances []struct {
		Employee  int    `json:"employee"`
		Date      string `json:"date"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		Break     int    `json:"break"`
	} `json:"attendances"`
}

type PersonioAuthRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type PersonioAuthResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Token   string `json:"token"`
		Expires int    `json:"expires_in"`
		Scope   string `json:"scope"`
	} `json:"data"`
}

func NewPersonioServer() *ServerData {
	s := &ServerData{}
	return s

}

func (s *ServerData) Run(cfg config.EnvConfig) error {
	if s.StartTime.IsZero() {
		s.StartTime = time.Now()
		s.Success = true
		ulog.Infof("Set StartTime: %s", s.StartTime.Format(time.RFC822))
	}

	for {
		if !s.CustomTimeLock {
			s.StopTime = time.Now()
			ulog.Tracef("Set StopTime: %s", s.StopTime.Format(time.RFC822))
			time.Sleep(cfg.StopTimeLoopInterval)
		}

	}
}

func (s *ServerData) setBreakStart(breakTime time.Time, override bool) error {

	// check if breakStart is alredy set (only one break is supported atm)
	if !s.BreakStart.IsZero() && !override {
		return fmt.Errorf("BreakStart is already set: %s. To overwrite use '--yes'", s.BreakStart.Format(time.RFC822))
	}

	// check if breakStart > start-time
	if breakTime.Before(s.StartTime) {
		return fmt.Errorf("BreakStart is before StartTime, this is not possible")
	}

	if override {
		oldBreakStart := s.BreakStart
		s.BreakStart = breakTime
		ulog.Infof("Overwritten BreakStart %s with %s", oldBreakStart, s.BreakStart)
		return nil
	}

	s.BreakStart = breakTime
	ulog.Infof("Successfully set BreakStart: %s", s.BreakStart)
	return nil
}

func (s *ServerData) setBreakEnd(breakTime time.Time, override bool) error {
	// Check if BreakStart is set
	if s.BreakStart.IsZero() {
		return fmt.Errorf("BreakStart is not set. Cannot set BreakEnd if BreakStart is not set yet")
	}

	// Check if now is after BreakStart
	if breakTime.Before(s.BreakStart) {
		return fmt.Errorf("BreakEnd (%s) is before BreakStart (%s)", breakTime, s.BreakStart)
	}

	if s.BreakEnd.IsZero() {
		s.BreakEnd = breakTime
		ulog.Infof("Successfully set BreakEnd: %s", s.BreakEnd)
		return nil
	}

	if override {
		oldBreakEnd := s.BreakEnd
		s.BreakEnd = breakTime
		ulog.Infof("Overwritten Break End %s with %s", oldBreakEnd, s.BreakEnd)
		return nil
	}

	return fmt.Errorf("BreakEnd is already set: %s. To overwrite use '--yes'", s.BreakEnd.Format(time.RFC822))
}

func (s *ServerData) sendToPersonio(cfg config.EnvConfig) error {

	dateFormat := "2006-01-02"
	timeFormat := "15:04"

	// check if Start and Stop Time is correct
	if s.StartTime.IsZero() {
		return fmt.Errorf("No StartTime set")
	}
	if s.StopTime.IsZero() {
		return fmt.Errorf("No StopTime set")
	}
	if s.StartTime.Format(dateFormat) != s.StopTime.Format(dateFormat) {
		return fmt.Errorf("StartTime and StopTime not on same day.")
	}
	if s.StartTime.After(s.StopTime) {
		return fmt.Errorf("StartTime after StopTime")
	}

	// check if break times are OK and get diff
	diff := time.Duration(0)
	if !s.BreakStart.IsZero() && !s.BreakEnd.IsZero() {

		if s.BreakEnd.After(s.StopTime) {
			return fmt.Errorf("BreakEnd is after of end-time")
		}
		if s.StartTime.Format(dateFormat) != s.BreakStart.Format(dateFormat) {
			return fmt.Errorf("StartTime and BreakStart not on same day")
		}
		if s.StartTime.Format(dateFormat) != s.BreakEnd.Format(dateFormat) {
			return fmt.Errorf("StartTime and BreakEnd not on same day")
		}

		diff = s.BreakEnd.Sub(s.BreakStart)
	}

	body := PersonioRequest{
		Attendances: []struct {
			Employee  int    `json:"employee"`
			Date      string `json:"date"`
			StartTime string `json:"start_time"`
			EndTime   string `json:"end_time"`
			Break     int    `json:"break"`
		}{
			{
				Employee:  cfg.PersonioEmployeeId,
				Date:      s.StartTime.Format(dateFormat),
				StartTime: s.StartTime.Format(timeFormat),
				EndTime:   s.StopTime.Format(timeFormat),
				Break:     int(diff.Minutes()),
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	authBody := PersonioAuthRequest{
		ClientId:     cfg.PersonioClientId,
		ClientSecret: cfg.PersonioClientSecret,
	}

	jsonAuthBody, err := json.Marshal(authBody)
	if err != nil {
		return err
	}

	authReq, err := http.NewRequest("POST", cfg.PersonioHost+"/v1/auth", bytes.NewBuffer(jsonAuthBody))
	if err != nil {
		return err
	}

	authReq.Header.Set("accept", "application/json")
	authReq.Header.Set("content-type", "application/json")

	authRes, err := http.DefaultClient.Do(authReq)
	if err != nil {
		return err
	}
	defer authRes.Body.Close()

	if authRes.StatusCode != 200 {
		return fmt.Errorf(authRes.Status)
	}

	authResBody, err := io.ReadAll(authRes.Body)
	if err != nil {
		return err
	}
	var authResData PersonioAuthResponse
	if err := json.Unmarshal(authResBody, &authResData); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", cfg.PersonioHost+"/v1/company/attendances", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResData.Data.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		resBody, _ := io.ReadAll(res.Body)
		ulog.Warnf(string(resBody))
		return fmt.Errorf(res.Status)
	}

	return nil

}
