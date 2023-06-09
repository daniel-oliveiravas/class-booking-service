package bookings

import (
	"time"
)

type Booking struct {
	ID        string    `json:"ID,omitempty"`
	MemberID  string    `json:"memberID,omitempty"`
	ClassID   string    `json:"classID,omitempty"`
	ClassDate time.Time `json:"classDate"`
	BookedAt  time.Time `json:"bookedAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type BookClass struct {
	MemberID  string    `json:"memberID,omitempty"`
	ClassID   string    `json:"classID,omitempty"`
	ClassDate time.Time `json:"classDate,omitempty"`
}

type PageInfo struct {
	Limit int
	Page  int
}
