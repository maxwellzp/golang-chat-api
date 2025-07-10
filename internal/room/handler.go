package room

import (
	"encoding/json"
	"github.com/maxwellzp/golang-chat-api/internal/httpx"
	"github.com/maxwellzp/golang-chat-api/internal/logger"
	"github.com/maxwellzp/golang-chat-api/internal/validatorx"
	"net/http"
)

type RoomHandler struct {
	roomService *RoomService
	validator   *validatorx.Validator
	logger      *logger.Logger
}

func NewRoomHandler(roomService *RoomService, validator *validatorx.Validator, logger *logger.Logger) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
		validator:   validator,
		logger:      logger,
	}
}

func (h *RoomHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httpx.GetUserID(r.Context())
		if err != nil {
			h.logger.Warnw("Unauthorized request to create room")
			httpx.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		var req CreateRoomRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.logger.Warnw("Failed to decode CreateRoomRequest",
				"error", err,
			)
			httpx.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if err := h.validator.Validate(&req); err != nil {
			h.logger.Warnw("Validation failed for CreateRoomRequest",
				"error", err,
				"user_id", userID,
			)
			httpx.WriteValidationError(w, err)
			return
		}

		rm, err := h.roomService.Create(r.Context(), userID, req)
		if err != nil {
			h.logger.Errorw("Failed to create room",
				"error", err,
				"user_id", userID,
			)
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		h.logger.Infow("Room created",
			"message_id", rm.ID,
			"user_id", userID,
		)
		httpx.WriteJSON(w, http.StatusCreated, rm)
	}
}

func (h *RoomHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httpx.GetUserID(r.Context())
		if err != nil {
			h.logger.Warnw("Unauthorized request to update room")
			httpx.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		id, err := httpx.ParseInt64Param(r, "id")
		if err != nil {
			h.logger.Warnw("Invalid room ID for update",
				"error", err,
			)
			httpx.WriteError(w, http.StatusBadRequest, "Invalid RoomID")
			return
		}

		var req UpdateRoomRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.logger.Warnw("Invalid UpdateRoomRequest body",
				"error", err,
			)
			httpx.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if err := h.validator.Validate(&req); err != nil {
			h.logger.Warnw("Validation failed for UpdateRoomRequest",
				"error", err,
				"user_id", userID,
				"room_id", id,
			)
			httpx.WriteValidationError(w, err)
			return
		}

		if err := h.roomService.Update(r.Context(), id, userID, req); err != nil {
			h.logger.Errorw("Failed to update room",
				"error", err,
				"user_id", userID,
				"message_id", id,
			)
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		h.logger.Infow("Room updated",
			"room_id", id,
			"user_id", userID,
		)
		httpx.WriteJSON(w, http.StatusNoContent, nil)
	}
}

func (h *RoomHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httpx.GetUserID(r.Context())
		if err != nil {
			h.logger.Warnw("Unauthorized request to delete room")
			httpx.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		id, err := httpx.ParseInt64Param(r, "id")
		if err != nil {
			h.logger.Warnw("Invalid room ID for delete",
				"error", err,
			)
			httpx.WriteError(w, http.StatusBadRequest, "Invalid MessageID")
			return
		}

		if err := h.roomService.Delete(r.Context(), id, userID); err != nil {
			h.logger.Errorw("Failed to delete room",
				"error", err,
				"user_id", userID,
				"room_id", id,
			)
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		h.logger.Infow("Room deleted",
			"room_id", id,
			"user_id", userID,
		)
		httpx.WriteJSON(w, http.StatusNoContent, nil)
	}
}

func (h *RoomHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.ParseInt64Param(r, "id")
		if err != nil {
			h.logger.Warnw("Invalid room ID for get",
				"error", err,
			)
			httpx.WriteError(w, http.StatusBadRequest, "Invalid MessageID")
			return
		}
		rm, err := h.roomService.GetByID(r.Context(), id)
		if err != nil {
			h.logger.Errorw("Failed to get room by ID",
				"error", err,
				"room_id", id,
			)
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		if rm == nil {
			h.logger.Warnw("Room not found",
				"room_id", id,
			)
			httpx.WriteError(w, http.StatusNotFound, "Room not found")
			return
		}
		h.logger.Infow("Room retrieved",
			"room_id", id,
		)
		httpx.WriteJSON(w, http.StatusOK, rm)
	}
}

func (h *RoomHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := h.roomService.List(r.Context())
		if err != nil {
			h.logger.Errorw("Failed to list room",
				"error", err,
			)
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		h.logger.Infow("Rooms listed",
			"count", len(rooms),
		)
		httpx.WriteJSON(w, http.StatusOK, rooms)
	}
}
