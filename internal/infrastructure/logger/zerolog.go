package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func InitLogger(logPath string) {
	if logPath == "" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		runLogFile, _ := os.OpenFile(
			logPath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		log.Logger = log.Output(runLogFile)
	}

	// Example of structured logging
	//log.Info().
	//	Str("event", "user_signup").
	//	Str("user", "johndoe").
	//	Time("timestamp", time.Now()).
	//	Msg("User signed up")
}
