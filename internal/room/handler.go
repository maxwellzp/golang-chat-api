package room

import (
	"fmt"
	"net/http"
)

type RoomHandler struct {
}

func NewRoomHandler(router *http.ServeMux) {
	handler := &RoomHandler{}
	router.HandleFunc("POST /room/create", handler.Create())
	router.HandleFunc("PATCH /room/update/{id}", handler.Update())
	router.HandleFunc("DELETE /room/delete/{id}", handler.Delete())
}

func (h *RoomHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Create a room")
	}
}

func (h *RoomHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Println("Update a room", id)
	}
}

func (h *RoomHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Println("Delete a room", id)
	}
}
