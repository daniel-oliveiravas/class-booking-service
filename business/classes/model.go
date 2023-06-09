package classes

import (
	"time"
)

type Class struct {
	ID        string    `json:"ID,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name,omitempty"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Capacity  int       `json:"capacity,omitempty"`
}

type NewClass struct {
	Name      string    `json:"name,omitempty"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Capacity  int       `json:"capacity,omitempty"`
}

type UpdateClass struct {
	Name       *string    `json:"name,omitempty"`
	StartDate  *time.Time `json:"startDate,omitempty"`
	EndDate    *time.Time `json:"endDate,omitempty"`
	Capability *int       `json:"capability,omitempty"`
}

type PageInfo struct {
	Limit int
	Page  int
}
