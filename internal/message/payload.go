package message

type CreateMessageRequest struct {
	SenderID   int    `json:"sender_id"`
	RoomID     *int   `json:"room_id,omitempty"`
	ReceiverID *int   `json:"receiver_id,omitempty"`
	Content    string `json:"content"`
}

type UpdateMessageRequest struct {
	MessageID int    `json:"message_id"`
	SenderID  int    `json:"sender_id"`
	Content   string `json:"content"`
}
