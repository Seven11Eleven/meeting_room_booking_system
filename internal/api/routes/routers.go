package routes

import (
	"net/http"
	"time"

	"github.com/Seven11Eleven/meeting_room_booking_system/internal/api/handlers"
	"github.com/Seven11Eleven/meeting_room_booking_system/internal/services"
	"github.com/Seven11Eleven/meeting_room_booking_system/internal/storage/postgresql"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(r *chi.Mux, timeout time.Duration, db *pgxpool.Pool) http.Handler {
	

	reservationStorage := postgresql.NewStorage(db)
	reservationService := services.NewReservationService(reservationStorage, timeout)

	handler := handlers.NewReservationHandler(reservationService)

	r.Route("/reservations", func(r chi.Router) {
		r.Get("/{room_id}", handler.GetReservationsByRoom)
		r.Post("/", handler.Reserve)
		r.Delete("/", handler.CancelReserve)
	})

	return r
}
