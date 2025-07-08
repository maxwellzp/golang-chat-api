package message

import "github.com/maxwellzp/golang-chat-api/internal/httpx"

type CreateMessageRequest struct {
	RoomID     *int64 `json:"room_id,omitempty"`
	ReceiverID *int64 `json:"receiver_id,omitempty"`
	Content    string `json:"content" validate:"required,min=3"`
}

func (r *CreateMessageRequest) Validate() error {
	if r.RoomID == nil && r.ReceiverID == nil {
		return httpx.ValidationErrorMap{
			"room_id":     "either room_id or receiver_id must be provided",
			"receiver_id": "either room_id or receiver_id must be provided",
		}
	}
	if r.RoomID != nil && r.ReceiverID != nil {
		return httpx.ValidationErrorMap{
			"room_id":     "only one of room_id or receiver_id must be provided",
			"receiver_id": "only one of room_id or receiver_id must be provided",
		}
	}

	return nil
}

type UpdateMessageRequest struct {
	Content string `json:"content" validate:"required,min=3"`
}
