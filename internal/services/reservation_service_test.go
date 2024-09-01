package services_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Seven11Eleven/meeting_room_booking_system/internal/config"
	"github.com/Seven11Eleven/meeting_room_booking_system/internal/domain/models"
	"github.com/Seven11Eleven/meeting_room_booking_system/internal/services"
	"github.com/Seven11Eleven/meeting_room_booking_system/internal/storage/postgresql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)




func TestReservationServiceCreate(t *testing.T) {
	ctx := context.Background()

	cfg := config.LoadTestConfig()

	db := postgresql.NewPool(cfg)
	defer db.Close()
	storage := postgresql.NewStorage(db)
	service := services.NewReservationService(storage, 2*time.Second)

	t.Run("successful reservation", func(t *testing.T) {
		_, err := db.Exec(ctx, "DELETE FROM reservations")
		require.NoError(t, err)

		reservation := &models.Reservation{
			RoomID:    "411",
			StartTime: time.Now().Add(1 * time.Hour),
			EndTime:   time.Now().Add(2 * time.Hour),
		}

		err = service.Create(ctx, reservation)
		assert.NoError(t, err)
	})

	t.Run("reservation conflict", func(t *testing.T) {
		_, err := db.Exec(ctx, "DELETE FROM reservations")
		require.NoError(t, err)

		reservation := &models.Reservation{
			RoomID:    "412",
			StartTime: time.Now().Add(1 * time.Hour),
			EndTime:   time.Now().Add(2 * time.Hour),
		}

		err = service.Create(ctx, reservation)
		assert.NoError(t, err)

		err = service.Create(ctx, reservation)
		assert.ErrorIs(t, err, models.ErrRoomAlreadyReservated)
	})

	t.Run("reservation conflict case 2", func(t *testing.T) {
		_, err := db.Exec(ctx, "DELETE FROM reservations")
		require.NoError(t, err)

		reservation := &models.Reservation{
			RoomID:    "412",
			StartTime: time.Now().Add(1 * time.Hour),
			EndTime:   time.Now().Add(2 * time.Hour),
		}
		reservationSecond := &models.Reservation{
			RoomID:    "412",
			StartTime: time.Now().Add(1*time.Hour + 15*time.Minute),
			EndTime:   time.Now().Add(3 * time.Hour),
		}

		err = service.Create(ctx, reservation)
		assert.NoError(t, err)

		err = service.Create(ctx, reservationSecond)
		assert.ErrorIs(t, err, models.ErrRoomAlreadyReservated)
	})

	t.Run("time not provided", func(t *testing.T) {
		reservation := &models.Reservation{
			RoomID:    "413",
			StartTime: time.Time{},
			EndTime:   time.Time{},
		}

		err := service.Create(ctx, reservation)
		assert.ErrorIs(t, err, models.ErrTimeNotProvided)
	})

	t.Run("start time is in the past", func(t *testing.T) {
		reservation := &models.Reservation{
			RoomID:    "414",
			StartTime: time.Now().Add(-1 * time.Hour),
			EndTime:   time.Now().Add(1 * time.Hour),
		}

		err := service.Create(ctx, reservation)
		assert.ErrorIs(t, err, models.ErrPastTime)
	})

	t.Run("end time before start time", func(t *testing.T) {
		reservation := &models.Reservation{
			RoomID:    "415",
			StartTime: time.Now().Add(2 * time.Hour),
			EndTime:   time.Now().Add(1 * time.Hour),
		}

		err := service.Create(ctx, reservation)
		assert.ErrorIs(t, err, models.ErrEndTimeBeforeStartTime)
	})

	t.Run("reservation duration exceeds 24 hours", func(t *testing.T) {
		reservation := &models.Reservation{
			RoomID:    "416",
			StartTime: time.Now().Add(1 * time.Hour),
			EndTime:   time.Now().Add(25 * time.Hour),
		}

		err := service.Create(ctx, reservation)
		assert.ErrorIs(t, err, models.ErrReservationTimeExceedingLimit)
	})

	t.Run("end of month reservation", func(t *testing.T) {
		_, err := db.Exec(ctx, "DELETE FROM reservations")
		require.NoError(t, err)

		// Тест для перехода на конец февраля (в невисокосный год)
		startTime := time.Date(2028, 2, 28, 22, 0, 0, 0, time.UTC)
		endTime := time.Date(2028, 2, 29, 1, 0, 0, 0, time.UTC)

		reservation := &models.Reservation{
			RoomID:    "417",
			StartTime: startTime,
			EndTime:   endTime,
		}

		err = service.Create(ctx, reservation)
		assert.NoError(t, err)
	})

	t.Run("midnight transition reservation", func(t *testing.T) {
		_, err := db.Exec(ctx, "DELETE FROM reservations")
		require.NoError(t, err)

		startTime := time.Now().Truncate(24 * time.Hour).Add(23 * time.Hour)
		endTime := startTime.Add(2 * time.Hour)

		reservation := &models.Reservation{
			RoomID:    "418",
			StartTime: startTime,
			EndTime:   endTime,
		}

		err = service.Create(ctx, reservation)
		assert.NoError(t, err)
	})
}


func TestConcurrentReservations(t *testing.T) {
	ctx := context.Background()

	cfg := config.LoadTestConfig()

	db := postgresql.NewPool(cfg) 
	defer db.Close()
	storage := postgresql.NewStorage(db)
	service := services.NewReservationService(storage, 2*time.Second)

	roomID := "418"
	startTime := time.Now().Add(1 * time.Hour)
	endTime := startTime.Add(1 * time.Hour)

	_, err := db.Exec(ctx, "DELETE FROM reservations")
	require.NoError(t, err)

	numGoroutines := 150

	var successCount, failureCount int
	var wg sync.WaitGroup

	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()

			reservation := &models.Reservation{
				RoomID:    roomID,
				StartTime: startTime,
				EndTime:   endTime,
			}

			err := service.Create(ctx, reservation)

			if err == nil {
				successCount++
			} else if err == models.ErrRoomAlreadyReservated {
				failureCount++
			} else {
				t.Error(err)
			}
		}()
	}

	wg.Wait()

	// если успешных броней больше чем 1, то логично что где-то произошло пересечение
	if successCount > 1 {
		t.Logf("Успешные бронирования: %d,  Пересечения: %d", successCount, failureCount)
		t.Error("error reserves")
	} else {
		t.Logf("Только одно бронирование прошло успешно. Проверка на пересечения пройдена.")
	}
}

func TestConcurrentReservationsWithDifferentPayloads(t *testing.T) {
	ctx := context.Background()

	cfg := config.LoadTestConfig()

	db := postgresql.NewPool(cfg) 
	defer db.Close()
	storage := postgresql.NewStorage(db)
	service := services.NewReservationService(storage, 2*time.Second)

	_, err := db.Exec(ctx, "DELETE FROM reservations")
	require.NoError(t, err)

	numGoroutines := 200

	var successCount, failureCount int
	var wg sync.WaitGroup

	wg.Add(numGoroutines)

	rooms := []string{"411", "412", "413", "414", "415", "416", "417", "418", "419", "410"}
	startTime := time.Now().Add(1 * time.Hour)
	duration := 1 * time.Hour

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()

			roomID := rooms[i%len(rooms)] 
			shift := time.Duration(i%15) * 10 * time.Minute 

			reservation := &models.Reservation{
				RoomID:    roomID,
				StartTime: startTime.Add(shift),
				EndTime:   startTime.Add(shift).Add(duration),
			}

			err := service.Create(ctx, reservation)

			if err == nil {
				successCount++
			} else if err == models.ErrRoomAlreadyReservated {
				failureCount++
			} else {
				t.Error(err)
			}
		}(i)
	}

	wg.Wait()

	t.Logf("Успешные брони: %d,  Пересечения: %d", successCount, failureCount)

	// должно быть хотя бы одного успешное бронирование, иначе функция сама по себе не работает
	if successCount == 0 {
		t.Error("waiting at least 1 sucessful reservation")
	}

	//должно быть хотя бы одно пересечение, иначе это значит функция пропускает пересечения 
	if failureCount == 0 {
		t.Error("waiting at least 1 conflict ")
	}
}



func TestReservationServiceDelete(t *testing.T) {
	ctx := context.Background()

	cfg := config.LoadTestConfig()

	db := postgresql.NewPool(cfg)
	defer db.Close()
	storage := postgresql.NewStorage(db)
	service := services.NewReservationService(storage, 2*time.Second)

	t.Run("successful reservation deletion", func(t *testing.T) {
		_, err := db.Exec(ctx, "DELETE FROM reservations")
		require.NoError(t, err)

		reservation := &models.Reservation{
			RoomID:    "419",
			StartTime: time.Now().Add(1 * time.Hour),
			EndTime:   time.Now().Add(2 * time.Hour),
		}

		err = service.Create(ctx, reservation)
		assert.NoError(t, err)

		err = service.DeleteReservation(ctx, reservation)
		assert.NoError(t, err)
		roomReservations, err := service.GetByRoomID(ctx, reservation.RoomID)
		assert.NoError(t, err)
		assert.Empty(t, roomReservations.Reservations)
	})

	t.Run("deletion of non-existing reservation", func(t *testing.T) {
		_, err := db.Exec(ctx, "DELETE FROM reservations")
		require.NoError(t, err)

		reservation := &models.Reservation{
			RoomID:    "420",
			StartTime: time.Now().Add(1 * time.Hour),
			EndTime:   time.Now().Add(2 * time.Hour),
		}
		err = service.DeleteReservation(ctx, reservation)
		assert.ErrorIs(t, err, models.ErrNoMatchingReservation)
	})
}

func TestReservationServiceGetByRoomID(t *testing.T) {
	ctx := context.Background()

	cfg := config.LoadTestConfig()

	db := postgresql.NewPool(cfg)
	defer db.Close()
	storage := postgresql.NewStorage(db)
	service := services.NewReservationService(storage, 2*time.Second)

	t.Run("successful get by room ID", func(t *testing.T) {
		_, err := db.Exec(ctx, "DELETE FROM reservations")
		require.NoError(t, err)

		roomID := "421"
		startTime := time.Date(2024, time.September, 1, 10, 51, 5, 0, time.UTC)
		endTime := startTime.Add(1 * time.Hour)

		reservation := &models.Reservation{
			RoomID:    roomID,
			StartTime: startTime,
			EndTime:   endTime,
		}
		err = service.Create(ctx, reservation)
		assert.NoError(t, err)

		roomReservations, err := service.GetByRoomID(ctx, roomID)
		assert.NoError(t, err)
		assert.NotEmpty(t, roomReservations.Reservations)
		assert.Equal(t, roomID, roomReservations.RoomID)
		assert.Len(t, roomReservations.Reservations, 1)
		expectedStartTime := reservation.StartTime.UTC()
		expectedEndTime := reservation.EndTime.UTC()
		actualStartTime := roomReservations.Reservations[0].StartTime.UTC()
		actualEndTime := roomReservations.Reservations[0].EndTime.UTC()

		assert.Equal(t, expectedStartTime, actualStartTime)
		assert.Equal(t, expectedEndTime, actualEndTime)
	})

	t.Run("get by room ID with no reservations", func(t *testing.T) {
		_, err := db.Exec(ctx, "DELETE FROM reservations")
		require.NoError(t, err)

		roomID := "422"

		roomReservations, err := service.GetByRoomID(ctx, roomID)
		assert.NoError(t, err)
		assert.Empty(t, roomReservations.Reservations)
		assert.Equal(t, roomID, roomReservations.RoomID)
	})
}
