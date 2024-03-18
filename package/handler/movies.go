package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	filmoteka "vk_restAPI"
	logger "vk_restAPI/logs"
)

type CreateMoviSwaggerRequest struct {
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	ReleaseDate string `json:"release_date" db:"release_date"`
	Rating      int    `json:"rating" db:"rating"`
}

type MovieSwaggerRequest struct {
	Movie    CreateMoviSwaggerRequest `json:"movie"`
	ActorIDs []int                    `json:"actorIDs"`
}

type movieRequest struct {
	Movie    filmoteka.Movies `json:"movie"`
	ActorIDs []int            `json:"actorIDs"`
}

// @Summary Create Movie
// @Security ApiKeyAuth
// @Tags movies
// @Description Create a new movie
// @Accept json
// @Produce json
// @Param input body MovieSwaggerRequest true "Movie information"
// @Success 200 {string} string "id"
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /api/movies/create [post]
func (h *Handler) handleCreateMovie(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Create Movie request")

	if err := h.checkAdminStatus(w, r); err != nil {
		logger.Log.Error("Admin status is not available:", err.Error())
		NewErrorResponse(w, http.StatusForbidden, "This function is only available to the administrator")
		return
	}

	var request movieRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Log.Error("Failed to decaode request body:", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	actors, err := h.service.GetActors()
	if err != nil {
		logger.Log.Error("Failed to handle GetActors:", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
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
			logger.Log.Printf("actor ID not found: %v", actorID)
			continue
		}
	}

	id, err := h.service.Movies.CreateMovie(request.Movie, validateActorIDs)
	if err != nil {
		logger.Log.Error("Failed to create movie:", err.Error())
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

type getMoviesResponse struct {
	Data []filmoteka.MoviesWithActors `json:"data"`
}

// @Summary Get All Movies
// @Security ApiKeyAuth
// @Tags movies
// @Description Get List of Movies Sorted By Rating
// @Accept json
// @Produce json
// @Success 200 {object} getMoviesResponse
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /api/movies [get]
func (h *Handler) handleGetAllMovies(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Get All Movies")

	movies, err := h.service.MoviesWithActors.GetMovies()
	if err != nil {
		logger.Log.Error("Failed to Get All Movies: ", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(movies) == 0 {
		NewErrorResponse(w, http.StatusOK, "The list of movies is empty")
		return
	}

	response := getMoviesResponse{Data: movies}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}
}

// @Summary Get All Movies Sorted By Title
// @Security ApiKeyAuth
// @Tags movies
// @Description Get List of Movies Sorted By Title
// @Accept json
// @Produce json
// @Success 200 {object} getMoviesResponse
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /api/movies/sort/title [get]
func (h *Handler) handleGetAllMoviesSortedByTitle(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Get All Movies Sorted by Title")

	movies, err := h.service.MoviesWithActors.GetMoviesSortedByTitle()
	if err != nil {
		logger.Log.Error("Failed to Get All Movies Sorted by Title: ", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(movies) == 0 {
		NewErrorResponse(w, http.StatusOK, "The list of movies is empty")
		return
	}

	response := getMoviesResponse{Data: movies}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}
}

// @Summary Get All Movies Sorted By Date
// @Security ApiKeyAuth
// @Tags movies
// @Description Get List of Movies Sorted By Date
// @Accept json
// @Produce json
// @Success 200 {object} getMoviesResponse
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /api/movies/sort/date [get]
func (h *Handler) handleGetAllMoviesSortedByDate(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Get All Movies Sorted by Date")

	movies, err := h.service.MoviesWithActors.GetMoviesSortedByDate()
	if err != nil {
		logger.Log.Error("Failed to Get All Movies Sorted by Date: ", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(movies) == 0 {
		NewErrorResponse(w, http.StatusOK, "The list of movies is empty")
		return
	}

	response := getMoviesResponse{Data: movies}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}
}

// @Summary Get Movie By ID
// @Security ApiKeyAuth
// @Tags movies
// @Description Get Movie by ID
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} filmoteka.MoviesWithActors
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /api/movies/{id} [get]
func (h *Handler) handleGetMovieById(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Get Movie By ID")

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

	movie, err := h.service.MoviesWithActors.GetMovieById(id)
	if err != nil {
		logger.Log.Error("Failed to Get Movie By ID: ", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := movie
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}
}

// @Summary Update Movie
// @Security ApiKeyAuth
// @Tags movies
// @Description Update information about Movie
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Param input body filmoteka.UpdateMovies true "Movie information for update"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /api/movies/{id} [put]
func (h *Handler) handleUpdateMovie(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Update Movie")

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

	var input filmoteka.UpdateMovies
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Log.Error("Failed to decode request body: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.UpdateMovie(id, input); err != nil {
		logger.Log.Error("Failed to update movie: ", err.Error())
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

// @Summary Delete Movie by Id
// @Security ApiKeyAuth
// @Tags movies
// @Description Delete information about Movie
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /api/movies/{id} [delete]
func (h *Handler) handleDeleteMovie(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Delete Movie")

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

	err = h.service.Movies.DeleteMovie(id)
	if err != nil {
		logger.Log.Error("Failed to delete movie: ", err.Error())
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

type Search struct {
	Fragment string `json:"fragment"`
}

// @Summary Search Movie By Fragment of Title
// @Security ApiKeyAuth
// @Tags movies
// @Description Search Movie By Title
// @Accept json
// @Produce json
// @Param input body search true "Search Movie By Fragment Of Title"
// @Success 200 {object} getMoviesResponse
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /api/movies/searchbytitle [post]
func (h *Handler) handleSearchMoviesByTitle(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Search Movie By Title")

	var fragmentTitle Search
	if err := json.NewDecoder(r.Body).Decode(&fragmentTitle); err != nil {
		logger.Log.Error("Failed to decode request body: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	movies, err := h.service.SearchMoviesByTitle(string(fragmentTitle.Fragment))
	if err != nil {
		logger.Log.Error("Failed to search movie by title: ", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(movies) == 0 {
		NewErrorResponse(w, http.StatusOK, "The list of movies is empty")
		return
	}

	response := getMoviesResponse{Data: movies}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}
}

// @Summary Search Movie By Fragment Of Actor Name
// @Security ApiKeyAuth
// @Tags movies
// @Description Search Movie By Fragment Of Actor Name
// @Accept json
// @Produce json
// @Param input body search true "Search Movie By Fragment Of Actor Name"
// @Success 200 {object} getMoviesResponse
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /api/movies/searchbyactor [post]
func (h *Handler) handleSearchMoviesByActorName(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Search Movie By Actor Name")

	var fragmentTitle Search
	if err := json.NewDecoder(r.Body).Decode(&fragmentTitle); err != nil {
		logger.Log.Error("Failed to decode request body: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	movies, err := h.service.SearchMovieByActorName(string(fragmentTitle.Fragment))
	if err != nil {
		logger.Log.Error("Failed to search movie by actor name: ", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(movies) == 0 {
		NewErrorResponse(w, http.StatusOK, "The list of movies is empty")
		return
	}

	response := getMoviesResponse{Data: movies}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}
}
