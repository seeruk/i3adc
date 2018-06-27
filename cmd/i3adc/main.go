package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/seeruk/i3adc/daemon"
	"github.com/seeruk/i3adc/i3"
	"github.com/seeruk/i3adc/i3adc"
	"github.com/seeruk/i3adc/xrandr"
)

func main() {
	resolver := i3adc.NewResolver()

	logger := resolver.ResolveLogger()
	logger = logger.With("module", "main")
	logger.Info("i3adc starting...")

	ctx, cfn := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	backend := resolver.ResolveStateBackend()

	// TODO(seeruk): Should these be moved to the resolver? (Yes)
	i3Thread, i3EventCh := i3.NewThread(resolver.ResolveLogger())
	i3ThreadDone := daemon.NewBackgroundThread(ctx, i3Thread)

	xrandrThread := xrandr.NewThread(backend, resolver.ResolveLogger(), i3EventCh)
	xrandrThreadDone := daemon.NewBackgroundThread(ctx, xrandrThread)

	select {
	case sig := <-signals:
		fmt.Println() // Skip the ^C
		logger.Infow("stopping background threads", "signal", sig)
	case res := <-i3ThreadDone:
		logger.Fatalw("error starting i3 thread", "error", res.Error())
	case res := <-xrandrThreadDone:
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
	<-xrandrThreadDone

	logger.Info("i3adc exiting...")
}
