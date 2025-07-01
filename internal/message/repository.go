package message

import (
	"context"
	"github.com/maxwellzp/golang-chat-api/internal/db"
)

type MessageRepository struct {
	Database *db.Db
}

func NewMessageRepository(database *db.Db) *MessageRepository {
	return &MessageRepository{Database: database}
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
