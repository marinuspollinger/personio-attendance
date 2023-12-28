package pserver

import (
	"net/http"
	"personio-server/config"
	"time"

	"github.com/dunv/uhttp"
	ulog "github.com/dunv/ulog/v2"
)

type setBreakCustomTimeRequest struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

var GetStatusHandler = func(s *ServerData) uhttp.Handler {
	return uhttp.NewHandler(
		uhttp.WithGet(func(r *http.Request, returnCode *int) interface{} {
			s.Success = true
			s.Error = ""
			return s
		}),
	)
}

var StartBreakNowHandler = func(s *ServerData) uhttp.Handler {
	return uhttp.NewHandler(
		uhttp.WithGet(func(r *http.Request, returnCode *int) interface{} {
			s.Success = true
			s.Error = ""
			err := s.setBreakStart(time.Now())
			if err != nil {
				ulog.Errorf("Error setting BreakStart: %s", err.Error())
				s.Success = false
			}

			return s
		}),
	)
}

var EndBreakNowHandler = func(s *ServerData) uhttp.Handler {
	return uhttp.NewHandler(
		uhttp.WithGet(func(r *http.Request, returnCode *int) interface{} {
			s.Success = true
			s.Error = ""
			err := s.setBreakEnd(time.Now())
			if err != nil {
				s.Success = false
				s.Error = err.Error()
				ulog.Errorf("Error setting BreakEnd: %s", s.Error)
				return s
			}

			return s
		}),
	)
}

var SetBreakCustomTimeHandler = func(s *ServerData) uhttp.Handler {
	return uhttp.NewHandler(
		uhttp.WithPostModel(setBreakCustomTimeRequest{}, func(r *http.Request, model interface{}, returnCode *int) interface{} {
			s.Success = true
			s.Error = ""
			req := model.(*setBreakCustomTimeRequest)

			err := s.setBreakStart(req.StartTime)
			if err != nil {
				s.Success = false
				s.Error = err.Error()
				ulog.Errorf("Error setting BreakStart: %s", s.Error)
				return s
			}

			err = s.setBreakEnd(req.EndTime)
			if err != nil {
				s.Success = false
				s.Error = err.Error()
				ulog.Errorf("Error setting BreakEnd: %s", s.Error)
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
			s.CurrentTime = time.Now()

			return s
		}),
	)
}
