package models

import "errors"

var (
	// http status code - 409 Conflict
	ErrRoomAlreadyReservated = errors.New("this room for this time is already reservated")

	// http status code - 404 Not FOund
	ErrNoMatchingReservation = errors.New("no matching reservation found")

	// http status code - 400 Bad Request
	ErrPastTime                      = errors.New("provided time must be in future")
	ErrTimeNotProvided               = errors.New("start time and end time must be provded")
	ErrEndTimeBeforeStartTime        = errors.New("end time must be after start time")
	ErrReservationTimeExceedingLimit = errors.New("reservation duration cannot be more than 24 hours")
)
