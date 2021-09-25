package main

import (
	"context"
	"flag"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func app() error {
	appCtx, cancel := context.WithCancel(context.Background())

	var daemon bool
	flag.BoolVar(&daemon, "daemon", false, "daemon mode")
	flag.Parse()

	errChan := make(chan error)
	go func() {
		handleSIGINT()
		cancel()
	}()

	go func() {
		errChan <- run(cancel, daemon)
	}()
	select {
	case err := <-errChan:
		return err
	case <-appCtx.Done():
		return nil
	}
}




func run(cancel context.CancelFunc, daemon bool) error {

	if !daemon {
		//do something
		println("instance mode")
		cancel()
		time.Sleep(5 * time.Second)
	}

	println("daemon mode")
	return nil
}

func handleSIGUSR2(logger *logrus.Logger) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR2)
	for range ch {
		level := logger.GetLevel()
		switch level {
		case logrus.DebugLevel:
			logger.Warn("switching log level to INFO")
			logger.SetLevel(logrus.InfoLevel)
		default:
			logger.Warn("switching log level to DEBUG")
			logger.SetLevel(logrus.DebugLevel)
		}
	}
}

func handleSIGINT() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	for range sigCh {
		signal.Stop(sigCh)
		return
	}
}
