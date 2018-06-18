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
	"github.com/seeruk/i3adc/internal/output"
)

func main() {
	resolver := internal.NewResolver()

	logger := resolver.ResolveLogger()
	logger.Info("main: i3adc starting...")

	eventCh := make(chan struct{}, 1) // I wonder if this buffer should be larger...
	eventCh <- struct{}{}             // Always trigger a change at application startup.

	ctx, cfn := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	// TODO(seeruk): Should these be moved to the resolver?
	i3Thread := i3.NewThread(logger, eventCh)
	i3ThreadDone := daemon.NewBackgroundThread(ctx, i3Thread)

	outputThread := output.NewThread(logger, eventCh)
	outputThreadDone := daemon.NewBackgroundThread(ctx, outputThread)

	select {
	case sig := <-signals:
		fmt.Println() // Skip the ^C
		logger.Infow("stopping background threads", "signal", sig)
	case res := <-i3ThreadDone:
		logger.Fatalw("error starting i3 thread", "error", res.Error())
	case res := <-outputThreadDone:
		logger.Fatalw("error starting output thread", "error", res.Error())
	}

	cfn()

	go func() {
		time.AfterFunc(5*time.Second, func() {
			logger.Error("took too long stopping, exiting")
			os.Exit(1)
		})
	}()

	// Wait for our background threads to clean up.
	<-i3ThreadDone
	<-outputThreadDone

	logger.Info("main: i3adc exiting...")
}
