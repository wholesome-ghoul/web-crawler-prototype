package custom_logger

import "log"

func Log() *log.Logger {
	logger := log.New(
		log.Writer(),
		"",
		log.Ldate|log.Ltime,
	)

	return logger
}
