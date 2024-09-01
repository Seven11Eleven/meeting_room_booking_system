package services_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Seven11Eleven/meeting_room_booking_system/internal/config"
	"github.com/Seven11Eleven/meeting_room_booking_system/internal/domain/models"
	"github.com/Seven11Eleven/meeting_room_booking_system/internal/domain/models/mocks"
	"github.com/Seven11Eleven/meeting_room_booking_system/internal/services"
	"github.com/Seven11Eleven/meeting_room_booking_system/internal/storage/postgresql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)



func TestConcurrentReservations(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(mocks.ReservationStorage)
	service := services.NewReservationService(mockStorage, 2*time.Second)

	roomID := "411"
	startTime := time.Now().Add(1 * time.Hour)
	endTime := startTime.Add(1 * time.Hour)

	numGoroutines := 10

	t.Run("successful reservation", func(t *testing.T) {
		var wg sync.WaitGroup
		var mu sync.Mutex
		successCount := 0
		failureCount := 0

		wg.Add(numGoroutines)

		// Настройка моков для успешного и конфликтного бронирования
		mockStorage.On("IsReserved", mock.Anything, roomID, startTime, endTime).Return(false, nil).Once()
		mockStorage.On("IsReserved", mock.Anything, roomID, startTime, endTime).Return(true, nil).Times(numGoroutines - 1)
		mockStorage.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		// Запуск горутин для конкурентных бронирований
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				reservation := &models.Reservation{
					RoomID:    roomID,
					StartTime: startTime,
					EndTime:   endTime,
				}
				err := service.Create(ctx, reservation)
				mu.Lock()
				defer mu.Unlock()
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

		assert.Equal(t, 1, successCount, "expected 1 success")
		assert.Equal(t, numGoroutines-1, failureCount, "expected failures to equal numGoroutines - 1")
		mockStorage.AssertExpectations(t)
	})

	t.Run("all reservations conflict", func(t *testing.T) {
		var wg sync.WaitGroup
		var mu sync.Mutex
		successCount := 0
		failureCount := 0

		wg.Add(numGoroutines)

		// Настройка моков для всех конфликтных бронирований
		mockStorage.On("IsReserved", mock.Anything, roomID, startTime, endTime).Return(true, nil).Times(numGoroutines)

		// Запуск горутин для конкурентных бронирований
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				reservation := &models.Reservation{
					RoomID:    roomID,
					StartTime: startTime,
					EndTime:   endTime,
				}
				err := service.Create(ctx, reservation)
				mu.Lock()
				defer mu.Unlock()
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

		assert.Equal(t, 0, successCount, "expected 0 successes")
		assert.Equal(t, numGoroutines, failureCount, "expected all failures")
		mockStorage.AssertExpectations(t)
	})
}


func TestReservationServiceIntegration(t *testing.T) {
	ctx := context.Background()

	cfg := config.LoadTestConfig()

	db := postgresql.NewConn(cfg)
	defer db.Close(ctx)
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
	t.Run("reservation conflict case 2", func(t *testing.T) { //это в случае, если есть бронь к примеру на 13:00 - 14:00 и кто-то хочет сделать бронь на 13:15 - 15:00
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
}






func TestConcurrentReservationsWithoutMutex(t *testing.T) {
	ctx := context.Background()

	cfg := config.LoadTestConfig()

	db := postgresql.NewConn(cfg)
	defer db.Close(ctx)
	storage := postgresql.NewStorage(db)
	service := services.NewReservationService(storage, 2*time.Second)

	roomID := "418"
	startTime := time.Now().Add(1 * time.Hour)
	endTime := startTime.Add(1 * time.Hour)

	_, err := db.Exec(ctx, "DELETE FROM reservations")
	require.NoError(t, err)

	numGoroutines := 10

	var successCount, failureCount int
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()

			reservation := &models.Reservation{
				RoomID:    roomID,
				StartTime: startTime,
				EndTime:   endTime,
			}

			// Используем функцию CreateWithoutMutex вместо обычной Create
			err := service.Create(ctx, reservation)

			mu.Lock()
			if err == nil {
				successCount++
			} else if err == models.ErrRoomAlreadyReservated {
				failureCount++
			} else {
				t.Error(err)
			}
			mu.Unlock()
		}()
	}

	wg.Wait()

	// Ожидаем, что без мьютексов могут возникнуть пересечения, поэтому successCount может быть больше 1
	if successCount > 1 {
		t.Logf("Успешные бронирования: %d, Ожидаемые конфликты: %d", successCount, failureCount)
		t.Error("Пересекающиеся бронирования возможны без мьютексов!")
	} else {
		t.Logf("Только одно бронирование прошло успешно. Проверка на пересечения пройдена.")
	}
}