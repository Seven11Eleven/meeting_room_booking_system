package models

import (
	"context"
	"time"
)

type Reservation struct {
	ID        int       `json:"id"`
	RoomID    string    `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type TimeSlot struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type RoomReservations struct {
	RoomID       string     `json:"room_id"`
	Reservations []TimeSlot `json:"reservations"`
}

type ReservationService interface {
	Create(ctx context.Context, reservation *Reservation) error
	DeleteReservation(ctx context.Context, reservation *Reservation) error
	GetByRoomID(ctx context.Context, roomID string) (*RoomReservations, error)
}

type ReservationStorage interface {
	Create(ctx context.Context, reservation *Reservation) error
	DeleteReservation(ctx context.Context, reservation *Reservation) error
	GetByRoomID(ctx context.Context, roomID string) (*RoomReservations, error)
	IsReserved(ctx context.Context, roomID string, startTime, endTime time.Time) (bool, error)
}
