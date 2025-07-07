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

func (ms *MessageService) Create(ctx context.Context, userID int64, req CreateMessageRequest) (*Message, error) {
	if req.RoomID == nil && req.ReceiverID == nil {
		return nil, errors.New("either room_id or receiver_id must be provided")
	}

	msg := &Message{
		SenderID:   userID,
		RoomID:     req.RoomID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
	}

	if err := ms.messageRepository.Create(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (ms *MessageService) Update(ctx context.Context, id int64, userID int64, req UpdateMessageRequest) error {
	if req.Content == "" {
		return errors.New("content cannot be empty")
	}

	return ms.messageRepository.Update(ctx, id, userID, req.Content)
}

func (ms *MessageService) Delete(ctx context.Context, messageID, senderID int64) error {
	return ms.messageRepository.Delete(ctx, messageID, senderID)
}

func (ms *MessageService) GetByID(ctx context.Context, messageID int64, senderID int64) (*Message, error) {
	return ms.messageRepository.GetByID(ctx, messageID, senderID)
}

func (ms *MessageService) List(ctx context.Context, roomID *int64, receiverID *int64) ([]*Message, error) {
	if roomID == nil && receiverID == nil {
		return nil, errors.New("room_id or receiver_id must be specified")
	}
	return ms.messageRepository.List(ctx, roomID, receiverID)
}
