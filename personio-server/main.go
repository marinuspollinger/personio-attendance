package main

import (
	"personio-server/config"
	"personio-server/pserver"
	"time"

	"github.com/codingconcepts/env"
	"github.com/dunv/uhttp"
	ulog "github.com/dunv/ulog/v2"
)

func main() {
	cfg := config.EnvConfig{}
	ulog.FatalIfError(env.Set(&cfg))

	ulog.Configure(
		ulog.WithLogLevel(cfg.LogLevel),
		ulog.WithCallerFieldWidth(40),
		ulog.WithStripAdditionalFields(true),
	)

	ps := pserver.NewPersonioServer()
	go func() {
		ulog.FatalIfError(ps.Run(cfg))
	}()

	u := uhttp.NewUHTTP(
		uhttp.WithAddress(cfg.HttpAddress),
		uhttp.WithReadTimeout(15*time.Second),
		uhttp.WithReadHeaderTimeout(15*time.Second),
		uhttp.WithIdleTimeout(60*time.Second),
		uhttp.WithWriteTimeout(60*time.Second),
		uhttp.WithLogger(ulog.NewDefaultLogger()),
	)

	u.Handle("/api/getServerData", pserver.GetStatusHandler(ps))
	u.Handle("/api/sendToPersonio", pserver.SendToPersonioHandler(ps, cfg))

	u.Handle("/api/setBreakStartTime", pserver.SetBreakStartTimeHandler(ps))
	u.Handle("/api/setBreakStopTime", pserver.SetBreakStopTimeHandler(ps))
	u.Handle("/api/setStartTime", pserver.SetStartTimeHandler(ps))
	u.Handle("/api/setStopTime", pserver.SetStopTimeHandler(ps))

	ulog.FatalIfError(u.ListenAndServe())
}
