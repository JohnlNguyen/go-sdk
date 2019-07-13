package main

import (
	"log"
	"time"

	"go-sdk/logger"
	"go-sdk/web"
)

func main() {
	app := web.NewFromEnv().WithLogger(logger.All())

	app.GET("/", func(_ *web.Ctx) web.Result {
		return web.NoContent
	}, web.WithTimeout(500*time.Millisecond), web.JSONProviderAsDefault)

	app.GET("/for/:duration", func(r *web.Ctx) web.Result {
		duration, err := web.DurationValue(r.RouteParam("duration"))
		if err != nil {
			return web.JSON.BadRequest(err)
		}
		time.Sleep(duration)
		return web.NoContent
	}, web.WithTimeout(5*time.Second), web.JSONProviderAsDefault)

	app.GET("/panic", func(_ *web.Ctx) web.Result {
		panic("ONLY A TEST")
	}, web.WithTimeout(500*time.Millisecond), web.JSONProviderAsDefault)

	if err := web.StartWithGracefulShutdown(app); err != nil {
		log.Fatal(err)
	}
}
