package message

import (
	"context"
	"errors"
)

type MessageService struct {
	messageRepository *MessageRepository
}

func NewMessageService(messageRepository *MessageRepository) *MessageService {
	return &MessageService{messageRepository: messageRepository}
}

func (ms *MessageService) Create(ctx context.Context, req CreateMessageRequest) (*Message, error) {
	if req.RoomID == nil && req.ReceiverID == nil {
		return nil, errors.New("either room_id or receiver_id must be provided")
	}

	msg := &Message{
		SenderID:   req.SenderID,
		RoomID:     req.RoomID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
	}

	if err := ms.messageRepository.Create(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (ms *MessageService) Update(ctx context.Context, req UpdateMessageRequest) error {
	if req.Content == "" {
		return errors.New("content cannot be empty")
	}

	return ms.messageRepository.Update(ctx, req.MessageID, req.SenderID, req.Content)
}

func (ms *MessageService) Delete(ctx context.Context, messageID, senderID int) error {
	return ms.messageRepository.Delete(ctx, messageID, senderID)
}

func (ms *MessageService) GetByID(ctx context.Context, messageID int) (*Message, error) {
	return ms.messageRepository.GetByID(ctx, messageID)
}

func (ms *MessageService) List(ctx context.Context, roomID *int, receiverID *int) ([]*Message, error) {
	if roomID == nil && receiverID == nil {
		return nil, errors.New("room_id or receiver_id must be specified")
	}
	return ms.messageRepository.List(ctx, roomID, receiverID)
}
