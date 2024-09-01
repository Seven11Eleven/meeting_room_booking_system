package services

import (
	"context"
	"time"

	"github.com/Seven11Eleven/meeting_room_booking_system/internal/domain/models"
)

type reservationService struct {
	reservationStorage models.ReservationStorage
	contextTimeout     time.Duration
}

func TimeValidator(timeStart, timeEnd time.Time) error {
	if timeStart.IsZero() || timeEnd.IsZero() {
		return models.ErrTimeNotProvided
	}
	if timeStart.Before(time.Now()) {
		return models.ErrPastTime
	}
	if timeEnd.Before(timeStart) {
		return models.ErrEndTimeBeforeStartTime
	}
	if timeEnd.Sub(timeStart) > 24*time.Hour {
		return models.ErrReservationTimeExceedingLimit
	}
	return nil
}

// Create implements models.ReservationService.
func (r *reservationService) Create(ctx context.Context, reservation *models.Reservation) error {
	err := TimeValidator(reservation.StartTime, reservation.EndTime)
	if err != nil {
		return err
	}

	isReserved, err := r.reservationStorage.IsReserved(ctx, reservation.RoomID, reservation.StartTime, reservation.EndTime)
	if err != nil {
		return err
	}

	if isReserved {
		return models.ErrRoomAlreadyReservated
	}

	return r.reservationStorage.Create(ctx, reservation)
}

// DeleteReservation implements models.ReservationService.
func (r *reservationService) DeleteReservation(ctx context.Context, reservation *models.Reservation) error {
	if err := r.reservationStorage.DeleteReservation(ctx, reservation); err != nil{
		return err
	}
	return nil
}

// GetByRoomID implements models.ReservationService.
func (r *reservationService) GetByRoomID(ctx context.Context, roomID string) (*models.RoomReservations, error) {
	reservations, err := r.reservationStorage.GetByRoomID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	return reservations, nil
}



func NewReservationService(reservationStorage models.ReservationStorage, timeout time.Duration) models.ReservationService {
	return &reservationService{
		reservationStorage: reservationStorage,
		contextTimeout:     timeout,
	}
}





