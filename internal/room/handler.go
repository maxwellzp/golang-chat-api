package room

import (
	"encoding/json"
	"github.com/maxwellzp/golang-chat-api/internal/httpx"
	"net/http"
	"strconv"
)

type RoomHandler struct {
	roomService *RoomService
}

func NewRoomHandler(roomService *RoomService) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
	}
}

func (h *RoomHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httpx.GetUserID(r.Context())
		if err != nil {
			httpx.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		var req CreateRoomRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		rm, err := h.roomService.Create(r.Context(), userID, req)
		if err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, rm)
	}
}

func (h *RoomHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httpx.GetUserID(r.Context())
		if err != nil {
			httpx.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		rawID := r.PathValue("id")
		id, err := strconv.ParseInt(rawID, 10, 64)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid MessageID")
		}

		var req UpdateRoomRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if err := h.roomService.Update(r.Context(), id, userID, req); err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}

		httpx.WriteJSON(w, http.StatusNoContent, nil)
	}
}

func (h *RoomHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httpx.GetUserID(r.Context())
		if err != nil {
			httpx.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		rawID := r.PathValue("id")
		id, err := strconv.ParseInt(rawID, 10, 64)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid MessageID")
		}

		if err := h.roomService.Delete(r.Context(), id, userID); err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		httpx.WriteJSON(w, http.StatusNoContent, nil)
	}
}

func (h *RoomHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawID := r.PathValue("id")
		id, err := strconv.ParseInt(rawID, 10, 64)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid MessageID")
		}
		rm, err := h.roomService.GetByID(r.Context(), id)
		if err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		if rm == nil {
			httpx.WriteError(w, http.StatusNotFound, "Room not found")
			return
		}

		httpx.WriteJSON(w, http.StatusOK, rm)
	}
}

func (h *RoomHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := h.roomService.List(r.Context())
		if err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		httpx.WriteJSON(w, http.StatusOK, rooms)
	}
}
