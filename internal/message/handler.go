package message

import (
	"encoding/json"
	"fmt"
	"github.com/maxwellzp/golang-chat-api/internal/httpx"
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
			httpx.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		msg, err := h.messageService.Create(r.Context(), req)
		if err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, msg)
	}
}

func (h *MessageHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Println("Update a message", id)

		var req UpdateMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if err := h.messageService.Update(r.Context(), req); err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		httpx.WriteJSON(w, http.StatusNoContent, nil)
	}
}

func (h *MessageHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid message id")
			return
		}
		senderID := 1
		if err := h.messageService.Delete(r.Context(), id, senderID); err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		httpx.WriteJSON(w, http.StatusNoContent, nil)
	}
}

func (h *MessageHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		if idStr == "" {
			httpx.WriteError(w, http.StatusBadRequest, "Missing id parameter")
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid message id")
			return
		}
		msg, err := h.messageService.GetByID(r.Context(), id)
		if err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		if msg == nil {
			httpx.WriteError(w, http.StatusNotFound, "Message not found")
			return
		}
		httpx.WriteJSON(w, http.StatusOK, msg)
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
				httpx.WriteError(w, http.StatusBadRequest, "Invalid room_id")
				return
			}
			roomID = &id
		}

		if receiverIDStr != "" {
			id, err := strconv.Atoi(receiverIDStr)
			if err != nil {
				httpx.WriteError(w, http.StatusBadRequest, "Invalid receiver_id")
				return
			}
			receiverID = &id
		}

		messages, err := h.messageService.List(r.Context(), roomID, receiverID)
		if err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		httpx.WriteJSON(w, http.StatusOK, messages)
	}
}
