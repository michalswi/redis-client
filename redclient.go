package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/michalswi/keycloak_client/server"
	"github.com/michalswi/redis-client/apis"
	"github.com/michalswi/redis-client/rclient"
)

var version = "0.1.0"

func main() {
	fmt.Println("Go Redis Client")
	logger := log.New(os.Stdout, "redisClient ", log.LstdFlags|log.Lshortfile)

	Rhost := "localhost"
	Rport := "6379"
	ServiceAddr := os.Getenv("SERVICE_ADDR")
	APIPath := "/red"

	// redis client
	client := rclient.NewClient(Rhost, Rport)

	r := mux.NewRouter()
	prefix := r.PathPrefix(APIPath).Subrouter()
	srv := server.NewServer(prefix, ServiceAddr)

	prefix.Path("/home").Methods("GET").HandlerFunc(apis.Home(logger, version))
	prefix.Path("/ping").Methods("GET").HandlerFunc(apis.PingRedis(logger, client))
	prefix.Path("/setuser").Methods("POST").HandlerFunc(apis.SetRedis(logger, client))
	prefix.Path("/getuser/{uid}").Methods("GET").HandlerFunc(apis.GetRedis(logger, client))

	// start server
	go func() {
		logger.Printf("Starting server on port %s \n", ServiceAddr)
		err := srv.ListenAndServe()
		if err != nil {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	// shutdown server
	gracefulShutdown(srv, logger)
}

// graceful shutdown
func gracefulShutdown(srv *http.Server, logger *log.Logger) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	logger.Printf("Shutting down the server...\n")
	os.Exit(0)
}