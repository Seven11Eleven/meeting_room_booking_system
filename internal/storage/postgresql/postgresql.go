package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Seven11Eleven/meeting_room_booking_system/internal/config"
	"github.com/Seven11Eleven/meeting_room_booking_system/internal/domain/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

// DeleteReservation implements models.ReservationRepository.
func (s *Storage) DeleteReservation(ctx context.Context, reservation *models.Reservation) error {
	query := `
		WITH matched AS (
			SELECT 1
			FROM reservations
			WHERE room_id = $1
				AND start_time = $2
				AND end_time = $3
		)
		DELETE FROM reservations
		WHERE room_id = $1
			AND start_time = $2
			AND end_time = $3
			AND EXISTS (SELECT 1 FROM matched)
		RETURNING *;
	`

	res, err := s.db.Exec(ctx, query, reservation.RoomID, reservation.StartTime, reservation.EndTime)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNoMatchingReservation
	}

	return nil
}

// Create implements models.ReservationRepository.
func (s *Storage) Create(ctx context.Context, reservation *models.Reservation) error {
	query := `
		INSERT INTO reservations(room_id, start_time, end_time) VALUES($1, $2, $3)
	`

	_, err := s.db.Exec(ctx, query, reservation.RoomID, reservation.StartTime, reservation.EndTime)
	if err != nil {
		return err
	}
	return nil
}

// GetByRoomID implements models.ReservationRepository.
func (s *Storage) GetByRoomID(ctx context.Context, roomID string) (*models.RoomReservations, error) {
	query := `
		SELECT
		 		start_time, end_time
		FROM 
				reservations
		WHERE 
				room_id = $1
	`

	rows, err := s.db.Query(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservations := &models.RoomReservations{
		RoomID:       roomID,
		Reservations: []models.TimeSlot{},
	}
	for rows.Next() {
		var reservation models.TimeSlot
		err := rows.Scan(&reservation.StartTime, &reservation.EndTime)
		if err != nil {
			return nil, err
		}

		reservations.Reservations = append(reservations.Reservations, reservation)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

// IsReserved implements models.ReservationRepository.
func (s *Storage) IsReserved(ctx context.Context, roomID string, startTime time.Time, endTime time.Time) (bool, error) {
	query := `
	SELECT 
		COUNT(*)
	FROM 
		reservations
	WHERE
		room_id = $1
		AND 
		(
			(start_time < $3 AND end_time > $2)
		)
	`

	var cnt int
	err := s.db.QueryRow(ctx, query, roomID, startTime, endTime).Scan(&cnt)
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

func NewStorage(db *pgxpool.Pool) models.ReservationStorage {
	return &Storage{
		db: db,
	}
}
func NewPool(env *config.Config) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", env.DBUser, env.DBPass, env.DBHost, env.DBPort, env.DBName)

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}

	config.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	return pool
}

func Stop(pool *pgxpool.Pool) error {
	if pool == nil {
		return errors.New("not connected")
	}

	pool.Close()
	log.Println("connection pool to the database is closed")
	return nil
}
