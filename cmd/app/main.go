package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Seven11Eleven/meeting_room_booking_system/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	myApp, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("Failed to create the app: %v", err)
	}
	defer myApp.Close()

	go func() {
		if err := myApp.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to run the app: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down the app")
	stop()
	time.Sleep(1 * time.Second) 
}