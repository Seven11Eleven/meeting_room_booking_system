package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Seven11Eleven/meeting_room_booking_system/internal/api/routes"
	"github.com/Seven11Eleven/meeting_room_booking_system/internal/config"
	"github.com/Seven11Eleven/meeting_room_booking_system/internal/storage/postgresql"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type App struct {
	Router *chi.Mux
	DB     *pgx.Conn
	Env    *config.Config
}


func NewApp( ctx context.Context ) (*App, error){
	env := config.MustLoad()


	db := postgresql.NewConn(env)
	if db == nil{
		return nil, fmt.Errorf("failed to establish db conn")
	}

	router := chi.NewRouter()

	routes.SetupRoutes(router, 10*time.Second, db)
	

	return &App{
		Router: router,
		DB: db,
		Env: env,
	}, nil
}


func (a *App) Run() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", a.Env.AppPort),
		Handler:      a.Router, 
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Printf("Starting server on %s", server.Addr)
	return server.ListenAndServe()
}

// Close закрывает соединение с базой данных
func (a *App) Close() {
	if err := postgresql.Stop(a.DB); err != nil {
		log.Printf("Error closing the database connection: %v", err)
	}
	log.Println("Server and database connection closed")
}