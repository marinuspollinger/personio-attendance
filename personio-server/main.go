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

	u.Handle("/api/startBreakNow", pserver.StartBreakNowHandler(ps))
	u.Handle("/api/endBreakNow", pserver.EndBreakNowHandler(ps))

	u.Handle("/api/setCustomBreakStartTime", pserver.SetCustomBreakStartTime(ps))
	u.Handle("/api/setCustomBreakStopTime", pserver.SetCustomBreakStopTime(ps))
	u.Handle("/api/setCustomStartTime", pserver.SetCustomStartTime(ps))
	u.Handle("/api/setCustomStopTime", pserver.SetCustomStopTime(ps))

	ulog.FatalIfError(u.ListenAndServe())
}
