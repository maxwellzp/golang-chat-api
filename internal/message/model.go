package message

import (
	"time"
)

type Message struct {
	ID         int64     `json:"id"`
	SenderID   int       `json:"sender_id"`
	RoomID     *int      `json:"room_id,omitempty"`
	ReceiverID *int      `json:"receiver_id,omitempty"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
