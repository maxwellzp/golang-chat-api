package room

import "time"

type Room struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	IsPrivate bool      `json:"is_private"`
	CreatedBy *int64    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}
