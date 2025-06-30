package message

import "context"

type MessageRepository struct {
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{}
}

func (r *MessageRepository) Create(ctx context.Context, msg *Message) error {
	return nil
}

func (r *MessageRepository) Update(ctx context.Context, id int64) error {
	return nil
}

func (r *MessageRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func (r *MessageRepository) GetByID(ctx context.Context, id int64) (*Message, error) {
	return nil, nil
}

func (r *MessageRepository) List(ctx context.Context, limit, offset int) ([]*Message, error) {
	return nil, nil
}
