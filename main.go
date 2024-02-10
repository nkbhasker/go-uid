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

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	postgresUrl := os.Getenv("POSTGRES_URL")
	port := os.Getenv("PORT")

	idGenerator := NewIdGenerator()

	dbStore, err := Init(postgresUrl)
	defer func() {
		err := dbStore.CloseDB()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("DB connection closed")
	}()
	if err != nil {
		log.Fatal("Error connecting database")
	}
	userRepository := NewUserRepository(dbStore.DB(), idGenerator)
	userController := NewUserController(userRepository)

	router := InitRouter(userController)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	// Server run context
	srvCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancelShutdownCtx := context.WithTimeout(srvCtx, 30*time.Second)
		defer cancelShutdownCtx()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("Graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	log.Println(fmt.Sprintf("Server listening on port %s", port), "\nPress ctrl+c to stop")
	err = srv.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		log.Fatal("Server closed", err)
	}

	<-srvCtx.Done()
}
