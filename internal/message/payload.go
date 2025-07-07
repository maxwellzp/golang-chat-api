package message

type CreateMessageRequest struct {
	RoomID     *int64 `json:"room_id,omitempty"`
	ReceiverID *int64 `json:"receiver_id,omitempty"`
	Content    string `json:"content" validate:"required,min=3"`
}

type UpdateMessageRequest struct {
	Content string `json:"content" validate:"required,min=3"`
}
