package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	filmoteka "vk_restAPI"
)

func (h *Handler) handleCreateMovie(w http.ResponseWriter, r *http.Request) {
	if err := h.checkAdminStatus(w, r); err != nil {
		newErrorResponse(w, http.StatusForbidden, "This function is only available to the administrator")
		return
	}

	userId, err := getUserId(r)
	if err != nil {
		newErrorResponse(w, http.StatusUnauthorized, "user id not found")
		return
	}

	type movieRequest struct {
		Movie    filmoteka.Movies `json:"movie"`
		ActorIDs []int            `json:"actorIDs"`
	}

	var request movieRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	actors, err := h.service.GetActors()
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	mapActors := make(map[int]bool)
	for _, actor := range actors {
		mapActors[actor.Id] = true
	}

	validateActorIDs := make([]int, 0)
	for _, actorID := range request.ActorIDs {
		if mapActors[actorID] {
			validateActorIDs = append(validateActorIDs, actorID)
		} else {
			log.Printf("actor ID not found: %v", actorID)
			continue
		}
	}

	//TODO: ПРЕОБРАЗОВАТЬ STRING В TIME

	id, err := h.service.Movies.CreateMovie(userId, request.Movie, validateActorIDs)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"id": id,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

type getMoviesResponse struct {
	Data []filmoteka.MoviesWithActors `json:"data"`
}

func (h *Handler) handleGetAllMovies(w http.ResponseWriter, r *http.Request) {

	movies, err := h.service.MoviesWithActors.GetMovies()
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := getMoviesResponse{Data: movies}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleGetAllMoviesSortedByTitle(w http.ResponseWriter, r *http.Request) {

	movies, err := h.service.MoviesWithActors.GetMoviesSortedByTitle()
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := getMoviesResponse{Data: movies}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleGetAllMoviesSortedByDate(w http.ResponseWriter, r *http.Request) {

	movies, err := h.service.MoviesWithActors.GetMoviesSortedByDate()
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := getMoviesResponse{Data: movies}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleGetMovieById(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		newErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	idStr := parts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	movie, err := h.service.MoviesWithActors.GetMovieById(id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := movie
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleUpdateMovie(w http.ResponseWriter, r *http.Request) {
	if err := h.checkAdminStatus(w, r); err != nil {
		newErrorResponse(w, http.StatusForbidden, "This function is only available to the administrator")
		return
	}

	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		newErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	idStr := parts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	var input filmoteka.UpdateMovies
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.UpdateMovie(id, input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response := statusResponse{Status: "ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleDeleteMovie(w http.ResponseWriter, r *http.Request) {
	if err := h.checkAdminStatus(w, r); err != nil {
		newErrorResponse(w, http.StatusForbidden, "This function is only available to the administrator")
		return
	}

	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		newErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	idStr := parts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	err = h.service.Movies.DeleteMovie(id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := statusResponse{Status: "ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type search struct {
	Fragment string `json:"fragment"`
}

func (h *Handler) handleSearchMoviesByTitle(w http.ResponseWriter, r *http.Request) {
	var fragmentTitle search
	if err := json.NewDecoder(r.Body).Decode(&fragmentTitle); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	movies, err := h.service.SearchMoviesByTitle(string(fragmentTitle.Fragment))
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := getMoviesResponse{Data: movies}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleSearchMoviesByActorName(w http.ResponseWriter, r *http.Request) {
	var fragmentTitle search
	if err := json.NewDecoder(r.Body).Decode(&fragmentTitle); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	movies, err := h.service.SearchMovieByActorName(string(fragmentTitle.Fragment))
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := getMoviesResponse{Data: movies}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
