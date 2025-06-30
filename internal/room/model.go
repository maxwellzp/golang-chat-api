package room

type Room struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Private     bool   `json:"private"`
	Description string `json:"description"`
}
