// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	context "context"
	time "time"

	models "github.com/Seven11Eleven/meeting_room_booking_system/internal/domain/models"
	mock "github.com/stretchr/testify/mock"
)

// ReservationStorage is an autogenerated mock type for the ReservationStorage type
type ReservationStorage struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, reservation
func (_m *ReservationStorage) Create(ctx context.Context, reservation *models.Reservation) error {
	ret := _m.Called(ctx, reservation)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Reservation) error); ok {
		r0 = rf(ctx, reservation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteReservation provides a mock function with given fields: ctx, reservation
func (_m *ReservationStorage) DeleteReservation(ctx context.Context, reservation *models.Reservation) error {
	ret := _m.Called(ctx, reservation)

	if len(ret) == 0 {
		panic("no return value specified for DeleteReservation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Reservation) error); ok {
		r0 = rf(ctx, reservation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByRoomID provides a mock function with given fields: ctx, roomID
func (_m *ReservationStorage) GetByRoomID(ctx context.Context, roomID string) (*models.RoomReservations, error) {
	ret := _m.Called(ctx, roomID)

	if len(ret) == 0 {
		panic("no return value specified for GetByRoomID")
	}

	var r0 *models.RoomReservations
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.RoomReservations, error)); ok {
		return rf(ctx, roomID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.RoomReservations); ok {
		r0 = rf(ctx, roomID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.RoomReservations)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, roomID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsReserved provides a mock function with given fields: ctx, roomID, startTime, endTime
func (_m *ReservationStorage) IsReserved(ctx context.Context, roomID string, startTime time.Time, endTime time.Time) (bool, error) {
	ret := _m.Called(ctx, roomID, startTime, endTime)

	if len(ret) == 0 {
		panic("no return value specified for IsReserved")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time, time.Time) (bool, error)); ok {
		return rf(ctx, roomID, startTime, endTime)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time, time.Time) bool); ok {
		r0 = rf(ctx, roomID, startTime, endTime)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, time.Time, time.Time) error); ok {
		r1 = rf(ctx, roomID, startTime, endTime)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewReservationStorage creates a new instance of ReservationStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReservationStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *ReservationStorage {
	mock := &ReservationStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
