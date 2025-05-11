package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vrs-api/internal/config"
)

func Run() {
	// init gin router
	viperConfig := config.NewViper()
	router := config.NewGin()

	
	// Create http.Server
	port := ":"+viperConfig.GetString("SERVER_ADDRESS")
	server := &http.Server{
		Addr:    port,
		Handler: router.Handler(),
	}
	go func() {
		// service connections
		log.Printf("Server started, listen on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown Server ...")
	timeout := viperConfig.GetInt("GRACEFUL_TIMEOUT")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Printf("timeout of %d seconds. \n", timeout)
	log.Println("Server exiting")
}
