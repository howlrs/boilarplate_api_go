package models

import "time"

const (
	COLLECTION_RESERVATION = "reservations"
)

type ReserveState int

const (
	ReserveStateUndefined ReserveState = iota
	ReserveStateReserve
	ReserveStateReservating
	ReserveStateCompleted
)
const (
	ReserveStateCancel ReserveState = -(iota + 1)
	ReserveStateCanceling
	ReserveStateCanceled
	ReserveStateFailed
)

type ReservatedTime struct {
	ID string `json:"id" firestore:"id"`

	Status ReserveState `json:"status" firestore:"status"`

	Content    string `json:"content" firestore:"content"`
	SubContent string `json:"sub_content" firestore:"sub_content"`

	StartTime time.Time `json:"start_time" firestore:"start_time"`
	EndTime   time.Time `json:"end_time" firestore:"end_time"`
	Timezone  string    `json:"timezone" firestore:"timezone"`
}

func NewReservatedTime(startTime time.Time, endTime time.Time, timezone string) *ReservatedTime {
	return &ReservatedTime{
		Status:    ReserveStateUndefined,
		StartTime: startTime,
		EndTime:   endTime,
		Timezone:  timezone,
	}
}

func (r *ReservatedTime) ToCollection(isTest bool) string {
	if isTest {
		return "test_" + COLLECTION_RESERVATION
	}
	return COLLECTION_RESERVATION
}

func (r *ReservatedTime) Reserve() {
	r.Status = ReserveStateReserve
}

func (r *ReservatedTime) Reservating() {
	r.Status = ReserveStateReservating
}

func (r *ReservatedTime) Completed() {
	r.Status = ReserveStateCompleted
}

func (r *ReservatedTime) Cancel() {
	r.Status = ReserveStateCancel
}

func (r *ReservatedTime) Canceling() {
	r.Status = ReserveStateCanceling
}

func (r *ReservatedTime) Canceled() {
	r.Status = ReserveStateCanceled
}

func (r *ReservatedTime) Failed() {
	r.Status = ReserveStateFailed
}

func (r *ReservatedTime) IsReserved() bool {
	return r.Status > ReserveStateUndefined
}
