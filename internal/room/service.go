package room

import (
	"context"
)

type RoomService struct {
	roomRepository *RoomRepository
}

func NewRoomService(roomRepository *RoomRepository) *RoomService {
	return &RoomService{roomRepository: roomRepository}
}

func (rs *RoomService) Create(ctx context.Context, userId int64, req CreateRoomRequest) (*Room, error) {
	rm := &Room{
		Name:      req.Name,
		CreatedBy: &userId,
		IsPrivate: req.Private,
	}

	if err := rs.roomRepository.Create(ctx, rm); err != nil {
		return nil, err
	}
	return rm, nil
}

func (rs *RoomService) Update(ctx context.Context, roomID int64, userID int64, req UpdateRoomRequest) error {

	return rs.roomRepository.Update(ctx, roomID, userID, req.Name, req.Private)
}

func (rs *RoomService) Delete(ctx context.Context, roomID int64, userID int64) error {
	return rs.roomRepository.Delete(ctx, roomID, userID)
}

func (rs *RoomService) GetByID(ctx context.Context, roomID int64) (*Room, error) {
	return rs.roomRepository.GetByID(ctx, roomID)
}

func (rs *RoomService) List(ctx context.Context) ([]*Room, error) {
	return rs.roomRepository.List(ctx)
}
