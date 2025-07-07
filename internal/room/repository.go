package room

import (
	"context"
	"database/sql"
	"errors"
	"github.com/maxwellzp/golang-chat-api/internal/db"
	"time"
)

type RoomRepository struct {
	database *db.Db
}

func NewRoomRepository(database *db.Db) *RoomRepository {
	return &RoomRepository{database: database}
}

func (r *RoomRepository) Create(ctx context.Context, room *Room) error {
	query := `INSERT INTO rooms (name, is_private, created_by, created_at) 
				VALUES ($1, $2, $3, $4)
				RETURNING id, created_at;`

	row := r.database.QueryRowContext(ctx, query, room.Name, room.IsPrivate, room.CreatedBy, time.Now())
	err := row.Scan(&room.ID, &room.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoomRepository) Update(ctx context.Context, id int64, userID int64, name string, isPrivate bool) error {
	query := `
		UPDATE rooms SET name = $1, is_private = $2 WHERE id = $3 AND created_by = $4;
`
	res, err := r.database.ExecContext(ctx, query, name, isPrivate, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no room found or permission denied")
	}

	return nil
}

func (r *RoomRepository) Delete(ctx context.Context, roomID int64, userID int64) error {
	res, err := r.database.ExecContext(ctx,
		"DELETE FROM rooms WHERE id = $1 AND created_by = $2", roomID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no room found or permission denied")
	}
	return nil
}

func (r *RoomRepository) GetByID(ctx context.Context, roomID int64) (*Room, error) {
	query := `SELECT id, name, is_private, created_by, created_at FROM rooms WHERE id = $1`
	row := r.database.QueryRowContext(ctx, query, roomID)

	var rm Room
	err := row.Scan(&rm.ID, &rm.Name, &rm.IsPrivate, &rm.CreatedBy, &rm.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rm, nil
}

func (r *RoomRepository) List(ctx context.Context) ([]*Room, error) {
	query := `SELECT id, name, is_private, created_by, created_at FROM rooms
`
	rows, err := r.database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*Room
	for rows.Next() {
		var msg Room
		err := rows.Scan(
			&msg.ID,
			&msg.Name,
			&msg.IsPrivate,
			&msg.CreatedBy,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, &msg)
	}
	return rooms, nil
}
