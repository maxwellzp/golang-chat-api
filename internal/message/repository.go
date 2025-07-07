package message

import (
	"context"
	"database/sql"
	"errors"
	"github.com/maxwellzp/golang-chat-api/internal/db"
	"time"
)

type MessageRepository struct {
	database *db.Db
}

func NewMessageRepository(database *db.Db) *MessageRepository {
	return &MessageRepository{database: database}
}

func (r *MessageRepository) Create(ctx context.Context, msg *Message) error {
	query := `
			INSERT INTO messages (sender_id, room_id, receiver_id, content, created_at, updated_at) 
			VALUES($1, $2, $3, $4, $5, $6) 
			RETURNING id, created_at, updated_at;`
	err := r.database.QueryRowContext(ctx, query,
		msg.SenderID,
		msg.RoomID,
		msg.ReceiverID,
		msg.Content,
		time.Now(),
		time.Now()).Scan(&msg.ID, &msg.CreatedAt, &msg.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *MessageRepository) Update(ctx context.Context, messageID int64, senderID int64, content string) error {
	query := `
		UPDATE messages SET content = $1, updated_at = $2 WHERE id = $3 AND sender_id = $4;
`
	res, err := r.database.ExecContext(ctx, query, content, time.Now(), messageID, senderID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no message found or permission denied")
	}

	return nil
}

func (r *MessageRepository) Delete(ctx context.Context, messageID int64, senderID int64) error {
	res, err := r.database.ExecContext(ctx,
		"DELETE FROM messages WHERE id = $1 AND sender_id = $2;", messageID, senderID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no message found or permission denied")
	}
	return nil
}

func (r *MessageRepository) GetByID(ctx context.Context, messageID int64, senderID int64) (*Message, error) {
	query := `SELECT id, sender_id, room_id, receiver_id, content, created_at, updated_at 
			  FROM messages WHERE id = $1 AND sender_id = $2;`
	row := r.database.QueryRowContext(ctx, query, messageID, senderID)

	var msg Message
	err := row.Scan(&msg.ID, &msg.SenderID, &msg.RoomID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt, &msg.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &msg, nil
}

func (r *MessageRepository) List(ctx context.Context, roomID *int64, receiverID *int64) ([]*Message, error) {
	query := `SELECT id, sender_id, room_id, receiver_id, content, created_at, updated_at 
				FROM messages
			WHERE ($1::int IS NULL OR room_id = $1)
          	AND ($2::int IS NULL OR receiver_id = $2)
        	ORDER BY created_at ASC
`
	rows, err := r.database.QueryContext(ctx, query, roomID, receiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(
			&msg.ID,
			&msg.SenderID,
			&msg.RoomID,
			&msg.ReceiverID,
			&msg.Content,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}
