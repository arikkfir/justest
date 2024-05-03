package justest

import (
	"os"
	"os/signal"
)

var (
	interruptionSignal os.Signal = nil
)

func init() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	go func() {
		for s := range signalChan {
			interruptionSignal = s
			break
		}
	}()
}

//go:noinline
func verifyNotInterrupted(t T) {
	GetHelper(t).Helper()
	if interruptionSignal != nil {
		t.Fatalf("Process has been canceled via signal %s", interruptionSignal)
	}
}
