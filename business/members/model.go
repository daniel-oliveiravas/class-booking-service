package members

import (
	"time"
)

type MemberType string

type Member struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	Name      string    `json:"name,omitempty"`
}

type NewMember struct {
	Name string `json:"name,omitempty"`
}

type UpdateMember struct {
	Name *string `json:"name,omitempty"`
}

type PageInfo struct {
	Limit int
	Page  int
}
