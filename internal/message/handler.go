package message

import (
	"fmt"
	"net/http"
)

type MessageHandler struct {
}

func NewMessageHandler(router *http.ServeMux) {
	handler := &MessageHandler{}
	router.HandleFunc("POST /message/create", handler.Create())
	router.HandleFunc("PATCH /message/update/{id}", handler.Update())
}

func (h *MessageHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Create a message")
	}
}

func (h *MessageHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Println("Update a message", id)
	}
}
