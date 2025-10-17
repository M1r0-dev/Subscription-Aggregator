// entry point of service
package main

import (
	"log"

	"github.com/M1r0-dev/Subscription-Aggregator/config"
	"github.com/M1r0-dev/Subscription-Aggregator/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Config error %w", err)
	}
	app.Run(cfg)
}
