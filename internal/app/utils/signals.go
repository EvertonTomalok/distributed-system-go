package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func MakeDoneSignal() chan os.Signal {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	return done
}
