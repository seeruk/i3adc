package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/seeruk/i3adc/internal"
	"github.com/seeruk/i3adc/internal/daemon"
	"github.com/seeruk/i3adc/internal/i3"
)

func main() {
	resolver := internal.NewResolver()

	logger := resolver.ResolveLogger()
	logger.Info("main: i3adc starting...")

	eventCh := make(chan struct{}, 1)
	eventCh <- struct{}{} // Always trigger a change at application startup.

	// TODO(seeruk): This is a mock output.Thread, needs to be implemented.
	go func() {
		for {
			select {
			case <-eventCh:
				fmt.Println("some event occurred")
			}
		}
	}()

	ctx, cfn := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	outputEventThread := i3.NewOutputEventThread(logger, eventCh)
	outputEventDone := daemon.NewBackgroundThread(ctx, outputEventThread)

	select {
	case sig := <-signals:
		fmt.Println() // Skip the ^C
		logger.Infow("stopping background threads", "signal", sig)
	case res := <-outputEventDone:
		logger.Fatalw("error starting output event thread", "error", res.Error())
	}

	cfn()

	go func() {
		time.AfterFunc(5*time.Second, func() {
			logger.Error("took too long stopping, exiting")
			os.Exit(1)
		})
	}()

	// Wait for our background threads to clean up.
	<-outputEventDone

	logger.Info("main: i3adc exiting...")
}
