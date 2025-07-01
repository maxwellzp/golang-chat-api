package room

type RoomService struct {
	RoomRepository *RoomRepository
}

func NewRoomService(roomRepository *RoomRepository) *RoomService {
	return &RoomService{RoomRepository: roomRepository}
}
