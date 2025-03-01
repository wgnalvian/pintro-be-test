package exception

import (
	"log"
	"runtime/debug"

	"wgnalvian.com/payment-server/config"
)

func LogError(err error) {
	if config.LoadConfig().APP_DEBUG {
		log.Printf("Error: %v\n%s", err, debug.Stack())
	}
}
