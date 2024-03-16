package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	filmoteka "vk_restAPI"
	logger "vk_restAPI/logs"
)

type CreateActorRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
}

// @Summary Create Actor
// @Security ApiKeyAuth
// @Tags actors
// @Description Create a new actor
// @Accept json
// @Produce json
// @Param input body CreateActorRequest true "Actor information"
// @Success 200 {string} string "id"
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /api/actors/create [post]
func (h *Handler) handleCreateActor(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Create Actor request")

	if err := h.checkAdminStatus(w, r); err != nil {
		logger.Log.Error("Admin status is not available:", err.Error())
		NewErrorResponse(w, http.StatusForbidden, "This function is only available to the administrator")
		return
	}
	var input filmoteka.Actors
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Log.Error("Failed to decaode request body:", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.Actors.CreateActor(input)
	if err != nil {
		logger.Log.Error("Failed to create actor:", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"id": id,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}
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
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /api/actors [get]
func (h *Handler) handleGetAllActors(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Get All Actors")

	actors, err := h.service.ActorsWithMovies.GetActors()
	if err != nil {
		logger.Log.Error("Failed to Get All Actors: ", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := getActorsResponse{Data: actors}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}
}

// @Summary Get Actor By ID
// @Security ApiKeyAuth
// @Tags actors
// @Description Get Actor by ID
// @Accept json
// @Produce json
// @Param id path int true "Actor ID"
// @Success 200 {object} filmoteka.ActorsWithMovies
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /api/actors/{id} [get]
func (h *Handler) handleGetActorById(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Get Actors By ID")

	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		logger.Log.Error("Missing ID parameter")
		NewErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	idStr := parts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log.Error("Invailed ID parameter: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	actor, err := h.service.ActorsWithMovies.GetActorById(id)
	if err != nil {
		logger.Log.Error("Failed to Get Actor By ID: ", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := actor
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}
}

// @Summary Update Actor
// @Security ApiKeyAuth
// @Tags actors
// @Description Update information about Actor
// @Accept json
// @Produce json
// @Param id path int true "Actor ID"
// @Param input body filmoteka.UpdateActors true "Actor information for update"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /api/actors/{id} [put]
func (h *Handler) handleUpdateActor(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Update Actor")

	if err := h.checkAdminStatus(w, r); err != nil {
		logger.Log.Error("Admin status is not available:", err.Error())
		NewErrorResponse(w, http.StatusForbidden, "This function is only available to the administrator")
		return
	}

	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		logger.Log.Error("Missing ID parameter")
		NewErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	idStr := parts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log.Error("Invailid ID parameter: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	var input filmoteka.UpdateActors
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Log.Error("Failed to decode request body: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.UpdateActor(id, input); err != nil {
		logger.Log.Error("Failed to update actor: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response := StatusResponse{Status: "ok"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}
}

// @Summary Delete Actor by Id
// @Security ApiKeyAuth
// @Tags actors
// @Description Delete information about Actor
// @Accept json
// @Produce json
// @Param id path int true "Actor ID"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /api/actors/{id} [delete]
func (h *Handler) handleDeleteActor(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Delete Actor")

	if err := h.checkAdminStatus(w, r); err != nil {
		logger.Log.Error("Admin status is not available:", err.Error())
		NewErrorResponse(w, http.StatusForbidden, "This function is only available to the administrator")
		return
	}

	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		logger.Log.Error("Missing ID parameter")
		NewErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	idStr := parts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log.Error("Failed to decode request body: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	err = h.service.Actors.DeleteActor(id)
	if err != nil {
		logger.Log.Error("Failed to delete actor: ", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := StatusResponse{Status: "ok"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}
}
