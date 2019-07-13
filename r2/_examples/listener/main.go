package main

import (
	"fmt"
	"os"

	"go-sdk/logger"
	"go-sdk/r2"
)

func main() {
	log := logger.All()

	err := r2.New("https://google.com/robots.txt",
		r2.OptHeaderValue("X-Sent-By", "go-sdk/request2"),
		r2.OptCookieValue("r2-ray-id", "baileydog01"),
		r2.OptLogResponse(log),
	).Discard()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := log.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
