package main

import (
	"context"
	"errors"
	"flag"
	"interview-go/config"
	"interview-go/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var path string
	flag.StringVar(&path, "config", "./config.yaml", "filepath to a config file")
	flag.Parse()

	cfg, err := config.NewConfiguration(path)
	if err != nil {
		log.Fatalf("could not load configuration: %v", err)
	}

	sv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("new server: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(quit)

	errCh := make(chan error, 1)
	go func() {
		if serr := sv.StartServer(); serr != nil && !errors.Is(serr, http.ErrServerClosed) {
			errCh <- serr
		}
	}()

	select {
	case sig := <-quit:
		log.Printf("signal received: %v — shutting down...", sig)
	case e := <-errCh:
		log.Printf("server error: %v — shutting down...", e)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Printf("shutdown complete")
}
