package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	filmoteka "vk_restAPI"
)

// @Summary Create Actor
// @Security ApiKeyAuth
// @Tags actors
// @Description Create a new actor
// @Accept json
// @Produce json
// @Param input body filmoteka.Actors true "Actor information"
// @Success 200 {string} string "id"
// @Failure 400 {object} Err "Bad Request"
// @Failure 403 {object} Err "Forbidden"
// @Failure 500 {object} Err "Internal Server Error"
// @Router /api/actors/create [post]
func (h *Handler) handleCreateActor(w http.ResponseWriter, r *http.Request) {
	if err := h.checkAdminStatus(w, r); err != nil {
		NewErrorResponse(w, http.StatusForbidden, "This function is only available to the administrator")
		return
	}
	var input filmoteka.Actors
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.Actors.CreateActor(input)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"id": id,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type getActorsResponse struct {
	Data []filmoteka.ActorsWithMovies `json:"data"`
}

// @Summary Get All Actors
// @Security ApiKeyAuth
// @Tags actors
// @Description Get List of Actors
// @Accept json
// @Produce json
// @Success 200 {object} getActorsResponse
// @Failure 400 {object} Err "Bad Request"
// @Failure 403 {object} Err "Forbidden"
// @Failure 500 {object} Err "Internal Server Error"
// @Router /api/actors [get]
func (h *Handler) handleGetAllActors(w http.ResponseWriter, r *http.Request) {

	actors, err := h.service.ActorsWithMovies.GetActors()
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := getActorsResponse{Data: actors}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleGetActorById(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		NewErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	idStr := parts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	actor, err := h.service.ActorsWithMovies.GetActorById(id)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := actor
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleUpdateActor(w http.ResponseWriter, r *http.Request) {
	if err := h.checkAdminStatus(w, r); err != nil {
		NewErrorResponse(w, http.StatusForbidden, "This function is only available to the administrator")
		return
	}

	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		NewErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	idStr := parts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	var input filmoteka.UpdateActors
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.UpdateActor(id, input); err != nil {
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response := StatusResponse{Status: "ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleDeleteActor(w http.ResponseWriter, r *http.Request) {
	if err := h.checkAdminStatus(w, r); err != nil {
		NewErrorResponse(w, http.StatusForbidden, "This function is only available to the administrator")
		return
	}

	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		NewErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	idStr := parts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	err = h.service.Actors.DeleteActor(id)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := StatusResponse{Status: "ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
