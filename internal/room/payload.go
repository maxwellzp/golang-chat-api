package room

type CreateRoomRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=50"`
	Private bool   `json:"private"`
}

type UpdateRoomRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=50"`
	Private bool   `json:"private"`
}
