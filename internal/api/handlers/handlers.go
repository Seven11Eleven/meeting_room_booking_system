package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Seven11Eleven/meeting_room_booking_system/internal/domain/models"
	"github.com/go-chi/chi/v5"
)

type ReservationHandler struct {
	ReservationService models.ReservationService
}

func NewReservationHandler(service models.ReservationService) *ReservationHandler {
	return &ReservationHandler{
		ReservationService: service,
	}
}

func (h *ReservationHandler) Reserve(w http.ResponseWriter, r *http.Request) {
	var reservation models.Reservation
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err := h.ReservationService.Create(r.Context(), &reservation)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated) //201
}

func (h *ReservationHandler) GetReservationsByRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room_id")

	reservations, err := h.ReservationService.GetByRoomID(r.Context(), roomID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //200
	if err := json.NewEncoder(w).Encode(reservations); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ReservationHandler) CancelReserve(w http.ResponseWriter, r *http.Request) {
	var reservation models.Reservation
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err := h.ReservationService.DeleteReservation(r.Context(), &reservation)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent) //204
}

func (h *ReservationHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, models.ErrRoomAlreadyReservated):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, models.ErrNoMatchingReservation):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, models.ErrTimeNotProvided),
		errors.Is(err, models.ErrPastTime),
		errors.Is(err, models.ErrEndTimeBeforeStartTime), 
		errors.Is(err, models.ErrReservationTimeExceedingLimit):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

