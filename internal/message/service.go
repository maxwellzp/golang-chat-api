package message

type MessageService struct {
	MessageRepository *MessageRepository
}

func NewMessageService(messageRepository *MessageRepository) *MessageService {
	return &MessageService{MessageRepository: messageRepository}
}
