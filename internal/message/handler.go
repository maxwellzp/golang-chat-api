package message

import (
	"encoding/json"
	"github.com/maxwellzp/golang-chat-api/internal/httpx"
	"github.com/maxwellzp/golang-chat-api/internal/logger"
	"github.com/maxwellzp/golang-chat-api/internal/validatorx"
	"net/http"
	"strconv"
)

type MessageHandler struct {
	messageService *MessageService
	validator      *validatorx.Validator
	logger         *logger.Logger
}

func NewMessageHandler(messageService *MessageService, validator *validatorx.Validator, logger *logger.Logger) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		validator:      validator,
		logger:         logger,
	}
}

func (h *MessageHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httpx.GetUserID(r.Context())
		if err != nil {
			h.logger.Warnw("Unauthorized request to create message")
			httpx.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		var req CreateMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.logger.Warnw("Failed to decode CreateMessageRequest",
				"error", err,
			)
			httpx.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if err := h.validator.Validate(&req); err != nil {
			h.logger.Warnw("Validation failed for CreateMessageRequest",
				"error", err,
				"user_id", userID,
			)
			httpx.WriteValidationError(w, err)
			return
		}

		if err := req.Validate(); err != nil {
			h.logger.Warnw("Custom validation failed for CreateMessageRequest",
				"error", err,
				"user_id", userID,
			)
			httpx.WriteValidationError(w, err)
			return
		}

		msg, err := h.messageService.Create(r.Context(), userID, req)
		if err != nil {
			h.logger.Errorw("Failed to create message",
				"error", err,
				"user_id", userID,
			)
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		h.logger.Infow("Message created",
			"message_id", msg.ID,
			"user_id", userID,
		)
		httpx.WriteJSON(w, http.StatusCreated, msg)
	}
}

func (h *MessageHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httpx.GetUserID(r.Context())
		if err != nil {
			h.logger.Warnw("Unauthorized request to update message")
			httpx.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		id, err := httpx.ParseInt64Param(r, "id")
		if err != nil {
			h.logger.Warnw("Invalid message ID for update",
				"error", err,
			)
			httpx.WriteError(w, http.StatusBadRequest, "Invalid MessageID")
			return
		}

		var req UpdateMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.logger.Warnw("Invalid UpdateMessageRequest body",
				"error", err,
			)
			httpx.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if err := h.validator.Validate(&req); err != nil {
			h.logger.Warnw("Validation failed for UpdateMessageRequest",
				"error", err,
				"user_id", userID,
				"message_id", id,
			)
			httpx.WriteValidationError(w, err)
			return
		}

		if err := h.messageService.Update(r.Context(), id, userID, req); err != nil {
			h.logger.Errorw("Failed to update message",
				"error", err,
				"user_id", userID,
				"message_id", id,
			)
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		h.logger.Infow("Message updated",
			"message_id", id,
			"user_id", userID,
		)
		httpx.WriteJSON(w, http.StatusNoContent, nil)
	}
}

func (h *MessageHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httpx.GetUserID(r.Context())
		if err != nil {
			h.logger.Warnw("Unauthorized request to delete message")
			httpx.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		id, err := httpx.ParseInt64Param(r, "id")
		if err != nil {
			h.logger.Warnw("Invalid message ID for delete",
				"error", err,
			)
			httpx.WriteError(w, http.StatusBadRequest, "Invalid MessageID")
			return
		}
		if err := h.messageService.Delete(r.Context(), id, userID); err != nil {
			h.logger.Errorw("Failed to delete message",
				"error", err,
				"user_id", userID,
				"message_id", id,
			)
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		h.logger.Infow("Message deleted",
			"message_id", id,
			"user_id", userID,
		)
		httpx.WriteJSON(w, http.StatusNoContent, nil)
	}
}

func (h *MessageHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httpx.GetUserID(r.Context())
		if err != nil {
			h.logger.Warnw("Unauthorized request to get message")
			httpx.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		id, err := httpx.ParseInt64Param(r, "id")
		if err != nil {
			h.logger.Warnw("Invalid message ID for get",
				"error", err,
			)
			httpx.WriteError(w, http.StatusBadRequest, "Invalid MessageID")
			return
		}
		msg, err := h.messageService.GetByID(r.Context(), id, userID)
		if err != nil {
			h.logger.Errorw("Failed to get message by ID",
				"error", err,
				"user_id", userID,
				"message_id", id,
			)
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		if msg == nil {
			h.logger.Warnw("Message not found",
				"message_id", id,
				"user_id", userID,
			)
			httpx.WriteError(w, http.StatusNotFound, "Message not found")
			return
		}
		h.logger.Infow("Message retrieved",
			"message_id", id,
			"user_id", userID,
		)
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
				h.logger.Warnw("Invalid room_id in query param",
					"value", roomIDStr,
				)
				httpx.WriteError(w, http.StatusBadRequest, "Invalid room_id")
				return
			}
			roomID = &id
		}

		if receiverIDStr != "" {
			id, err := strconv.ParseInt(receiverIDStr, 10, 64)
			if err != nil {
				h.logger.Warnw("Invalid receiver_id in query param",
					"value", receiverIDStr,
				)
				httpx.WriteError(w, http.StatusBadRequest, "Invalid receiver_id")
				return
			}
			receiverID = &id
		}

		h.logger.Infow("Listing messages",
			"room_id", roomID,
			"receiver_id", receiverID,
		)
		messages, err := h.messageService.List(r.Context(), roomID, receiverID)
		if err != nil {
			h.logger.Errorw("Failed to list messages",
				"error", err,
			)
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		h.logger.Infow("Messages listed",
			"count", len(messages),
		)
		httpx.WriteJSON(w, http.StatusOK, messages)
	}
}
