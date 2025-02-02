package main

import (
	"algorhytm/language"
	english "algorhytm/language/english"
	"algorhytm/server"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	// Here we only pass os.Args and os.Getenv, but you could pass os.Stdin, etc., too.
	if err := run(ctx, os.Getenv); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, getenv func(string) string) error {
	// 1. Create a context that is canceled by SIGINT or SIGTERM
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 2. Pretend we read some config from environment variables
	config := server.Config{
		Host: getenv("APP_HOST"), // e.g., "0.0.0.0"
		Port: getenv("APP_PORT"), // e.g., "8080"
	}
	if config.Host == "" {
		config.Host = "0.0.0.0"
	}
	if config.Port == "" {
		config.Port = "8080"
	}
	properties, err := english.GetProperties()
	if err != nil {
		return err

	}
	English := language.Initialize(properties)

	// 3. Create whatever you need for your server
	//    e.g. logger, stores, Slack link store, etc.
	//    In this minimal example, we’ll just create the handler.
	srv := server.NewServer(&config, English)

	// 4. Set up the HTTP server
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}

	// 5. Start serving in a goroutine
	go func() {
		fmt.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// In a real app, you might handle the error more robustly.
			fmt.Fprintf(os.Stderr, "server error: %s\n", err)
		}
	}()

	// 6. Wait for the context to be cancelled, then gracefully shut down the server
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Block until we get a cancellation (Ctrl+C or SIGTERM).
		<-ctx.Done()

		// Give the server a 10-second deadline to shut down.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()

	// 7. Block until the shutdown goroutine completes
	wg.Wait()

	// No errors, so return nil
	return nil
}
