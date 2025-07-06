package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go-monolite/pkg/logger"
)

type Shutdownable interface {
	Shutdown(ctx context.Context) error
}

func GracefulShutdown(timeout time.Duration, components ...Shutdownable) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup
	for _, comp := range components {
		wg.Add(1)
		go func(c Shutdownable) {
			defer wg.Done()
			if err := c.Shutdown(ctx); err != nil {
				logger.Error(err, "shutdown error")
			} else {
				logger.Info("shutdown success")
			}
		}(comp)
	}

	wg.Wait()
}
