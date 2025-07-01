package room

import (
	"context"
	"github.com/maxwellzp/golang-chat-api/internal/db"
)

type RoomRepository struct {
	Database *db.Db
}

func NewRoomRepository(database *db.Db) *RoomRepository {
	return &RoomRepository{Database: database}
}

func (r *RoomRepository) Create(ctx context.Context, room *Room) error {
	return nil
}

func (r *RoomRepository) Update(ctx context.Context, id int64) error {
	return nil
}

func (r *RoomRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func (r *RoomRepository) GetByID(ctx context.Context, id int64) (*Room, error) {
	return nil, nil
}

func (r *RoomRepository) List(ctx context.Context, limit, offset int) ([]*Room, error) {
	return nil, nil
}
