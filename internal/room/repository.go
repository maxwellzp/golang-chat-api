package room

import "context"

type RoomRepository struct {
}

func NewRoomRepository() *RoomRepository {
	return &RoomRepository{}
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
