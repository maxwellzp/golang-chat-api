package room

type CreateRoomRequest struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

type UpdateRoomRequest struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
}
