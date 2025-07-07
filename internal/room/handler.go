package room

import (
	"encoding/json"
	"fmt"
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
		fmt.Println("Create a room")
		user := 1
		var req CreateRoomRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		rm, err := h.roomService.Create(r.Context(), user, req)
		if err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, rm)
	}
}

func (h *RoomHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid room id")
			return
		}

		var req UpdateRoomRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if err := h.roomService.Update(r.Context(), req, id); err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}

		httpx.WriteJSON(w, http.StatusNoContent, nil)
	}
}

func (h *RoomHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid room id")
			return
		}

		if err := h.roomService.Delete(r.Context(), id); err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "Something went wrong. Please try again later")
			return
		}
		httpx.WriteJSON(w, http.StatusNoContent, nil)
	}
}

func (h *RoomHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		if idStr == "" {
			httpx.WriteError(w, http.StatusBadRequest, "Missing room id parameter")
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid room id")
			return
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
