package message

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type MessageHandler struct {
	messageService *MessageService
}

func NewMessageHandler(router *http.ServeMux, messageService *MessageService) {
	handler := &MessageHandler{
		messageService: messageService,
	}
	router.HandleFunc("POST /messages/create", handler.Create())
	router.HandleFunc("PATCH /messages/update/{id}", handler.Update())
	router.HandleFunc("DELETE /messages/delete/{id}", handler.Delete())
	router.HandleFunc("GET /messages/{id}", handler.GetByID())
	router.HandleFunc("GET /messages/list", handler.List())
}

func (h *MessageHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		msg, err := h.messageService.Create(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(msg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *MessageHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Println("Update a message", id)

		var req UpdateMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := h.messageService.Update(r.Context(), req); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *MessageHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}
		senderID := 1
		if err := h.messageService.Delete(r.Context(), id, senderID); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *MessageHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		if idStr == "" {
			http.Error(w, "Missing id parameter", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}
		msg, err := h.messageService.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if msg == nil {
			http.Error(w, "Message not found", http.StatusNotFound)
		}

		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(msg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *MessageHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomIDStr := r.URL.Query().Get("room_id")
		receiverIDStr := r.URL.Query().Get("receiver_id")

		var roomID *int
		var receiverID *int

		if roomIDStr != "" {
			id, err := strconv.Atoi(roomIDStr)
			if err != nil {
				http.Error(w, "Invalid room_id", http.StatusBadRequest)
				return
			}
			roomID = &id
		}

		if receiverIDStr != "" {
			id, err := strconv.Atoi(receiverIDStr)
			if err != nil {
				http.Error(w, "Invalid receiver_id", http.StatusBadRequest)
				return
			}
			receiverID = &id
		}

		messages, err := h.messageService.List(r.Context(), roomID, receiverID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(messages); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
