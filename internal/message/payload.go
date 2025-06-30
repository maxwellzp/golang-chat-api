package message

type CreateMessageRequest struct {
	Message string `json:"message"`
}

type UpdateMessageRequest struct {
	Message string `json:"message"`
}
