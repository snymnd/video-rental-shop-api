package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"vrs-api/internal/constant"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Run() {
	// load env
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load env", err)
	}

	// init gin router
	r := gin.New()
	r.ContextWithFallback = true

	// start server
	port := os.Getenv(constant.ENV_SERVER_ADDRESS)

	// Create http.Server
	server := &http.Server{
		Addr:    port,
		Handler: r.Handler(),
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
	timeout, strconvErr := strconv.Atoi(os.Getenv(constant.ENV_GRACEFUL_TIMEOUT_KEY))
	const defaultTimeout = 1
	if strconvErr != nil {
		timeout = defaultTimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Printf("timeout of %d seconds. \n", timeout)
	log.Println("Server exiting")
}
