package pserver

import (
	"net/http"
	"personio-server/config"
	"time"

	"github.com/dunv/uhttp"
	ulog "github.com/dunv/ulog/v2"
)

type setCustomTimeRequest struct {
	BreakStartTime time.Time `json:"break_start_time"`
	BreakStopTime  time.Time `json:"break_end_time"`
	StartTime      time.Time `json:"start_time"`
	StopTime       time.Time `json:"end_time"`
	Override       bool      `json:"override"`
}

var GetStatusHandler = func(s *ServerData) uhttp.Handler {
	return uhttp.NewHandler(
		uhttp.WithGet(func(r *http.Request, returnCode *int) interface{} {
			return s
		}),
	)
}

var SetStartTimeHandler = func(s *ServerData) uhttp.Handler {
	return uhttp.NewHandler(
		uhttp.WithPostModel(setCustomTimeRequest{}, func(r *http.Request, model interface{}, returnCode *int) interface{} {
			s.Success = true
			s.Error = ""

			req := model.(*setCustomTimeRequest)
			s.StartTime = req.StartTime
			return s
		}),
	)
}

var SetStopTimeHandler = func(s *ServerData) uhttp.Handler {
	return uhttp.NewHandler(
		uhttp.WithPostModel(setCustomTimeRequest{}, func(r *http.Request, model interface{}, returnCode *int) interface{} {
			s.Success = true
			s.Error = ""

			req := model.(*setCustomTimeRequest)

			if req.StopTime.Before(s.StartTime) {
				s.Success = false
				s.Error = "Stop Time is before Start Time"
				ulog.Errorf("Error setting Stop Time: %s", s.Error)
				return s
			}

			s.CustomTimeLock = true
			s.StopTime = req.StopTime
			return s
		}),
	)
}

var SetBreakStartTimeHandler = func(s *ServerData) uhttp.Handler {
	return uhttp.NewHandler(
		uhttp.WithPostModel(setCustomTimeRequest{}, func(r *http.Request, model interface{}, returnCode *int) interface{} {
			s.Success = true
			s.Error = ""

			req := model.(*setCustomTimeRequest)

			err := s.setBreakStart(req.BreakStartTime, req.Override)
			if err != nil {
				s.Success = false
				s.Error = err.Error()
				ulog.Errorf("Error setting BreakStart: %s", s.Error)
				return s
			}

			return s
		}),
	)
}

var SetBreakStopTimeHandler = func(s *ServerData) uhttp.Handler {
	return uhttp.NewHandler(
		uhttp.WithPostModel(setCustomTimeRequest{}, func(r *http.Request, model interface{}, returnCode *int) interface{} {
			s.Success = true
			s.Error = ""

			req := model.(*setCustomTimeRequest)

			err := s.setBreakEnd(req.BreakStopTime, req.Override)
			if err != nil {
				s.Success = false
				s.Error = err.Error()
				ulog.Errorf("Resetting BreakStart, Error setting BreakEnd: %s", s.Error)
				return s
			}

			return s
		}),
	)
}

var SendToPersonioHandler = func(s *ServerData, cfg config.EnvConfig) uhttp.Handler {
	return uhttp.NewHandler(
		uhttp.WithGet(func(r *http.Request, returnCode *int) interface{} {
			s.Success = true
			s.Error = ""

			err := s.sendToPersonio(cfg)
			if err != nil {
				s.Success = false
				s.Error = err.Error()
				ulog.Errorf("Error sending Personio Request: %s", s.Error)

			}

			// Reset times, since we sent sucessfully to personio
			s.BreakEnd = time.Time{}
			s.BreakStart = time.Time{}
			s.StartTime = time.Now()
			s.StopTime = time.Now()
			s.CustomTimeLock = false

			return s
		}),
	)
}
