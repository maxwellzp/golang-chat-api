package message

import (
	"encoding/json"
	"github.com/maxwellzp/golang-chat-api/internal/httpx"
	"net/http"
	"strconv"
)

type MessageHandler struct {
	messageService *MessageService
}

func NewMessageHandler(messageService *MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

func (h *MessageHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httpx.GetUserID(r.Context())
		if err != nil {
			httpx.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		var req CreateMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		msg, err := h.messageService.Create(r.Context(), userID, req)
		if err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, msg)
	}
}

func (h *MessageHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httpx.GetUserID(r.Context())
		if err != nil {
			httpx.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		id, err := httpx.ParseInt64Param(r, "id")
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid MessageID")
			return
		}

		var req UpdateMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if err := h.messageService.Update(r.Context(), id, userID, req); err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		httpx.WriteJSON(w, http.StatusNoContent, nil)
	}
}

func (h *MessageHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httpx.GetUserID(r.Context())
		if err != nil {
			httpx.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		id, err := httpx.ParseInt64Param(r, "id")
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid MessageID")
			return
		}
		if err := h.messageService.Delete(r.Context(), id, userID); err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		httpx.WriteJSON(w, http.StatusNoContent, nil)
	}
}

func (h *MessageHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httpx.GetUserID(r.Context())
		if err != nil {
			httpx.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		id, err := httpx.ParseInt64Param(r, "id")
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid MessageID")
			return
		}
		msg, err := h.messageService.GetByID(r.Context(), id, userID)
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

		var roomID *int64
		var receiverID *int64

		if roomIDStr != "" {
			id, err := strconv.ParseInt(roomIDStr, 10, 64)
			if err != nil {
				httpx.WriteError(w, http.StatusBadRequest, "Invalid room_id")
				return
			}
			roomID = &id
		}

		if receiverIDStr != "" {
			id, err := strconv.ParseInt(receiverIDStr, 10, 64)
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
