package logger

import (
	"fmt"
	"log"

	"klubRanks/config"
)

func LogInfo(message string, args ...any) {
	log.Println("[INFO]", format(message, args...))
}

func LogError(message string, args ...any) {
	log.Println("[ERROR]", format(message, args...))
}

func LogDebug(message string, args ...any) {
	if config.AppConfig.Server.Log == "debug" {
		log.Println("[DEBUG]", format(message, args...))
	}
}

func format(message string, args ...any) string {
	if len(args) == 0 {
		return message
	}
	return message + " " + fmt.Sprint(args...)
}
