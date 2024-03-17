package handler

import (
	"net/http"
	"vk_restAPI/package/service"

	_ "vk_restAPI/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/swagger/", httpSwagger.Handler())

	auth := "/auth"

	mux.HandleFunc(auth+"/sign-up", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.handleSignUp(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(auth+"/log-in", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.handleSignIn(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	api := "/api"

	//Actors
	apiActors := api + "/actors"

	//POST for /api/actors/create
	mux.HandleFunc(apiActors+"/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.userIdentity(h.handleCreateActor)(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	//GET for /api/actors
	mux.HandleFunc(apiActors, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.userIdentity(h.handleGetAllActors)(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	//GET, PUT, DELETE for /api/actors/id
	mux.HandleFunc(apiActors+"/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.userIdentity(h.handleGetActorById)(w, r)
		} else if r.Method == http.MethodDelete {
			h.userIdentity(h.handleDeleteActor)(w, r)
		} else if r.Method == http.MethodPut {
			h.userIdentity(h.handleUpdateActor)(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	//Movies
	apiMovies := api + "/movies"

	//POST for /api/movies/create
	mux.HandleFunc(apiMovies+"/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.userIdentity(h.handleCreateMovie)(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	//GET for /api/movies
	mux.HandleFunc(apiMovies, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.userIdentity(h.handleGetAllMovies)(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	//GET for /api/movies/sort/title
	mux.HandleFunc(apiMovies+"/sort/title", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.userIdentity(h.handleGetAllMoviesSortedByTitle)(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	//GET for /api/movies/sort/date
	mux.HandleFunc(apiMovies+"/sort/date", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.userIdentity(h.handleGetAllMoviesSortedByDate)(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	//GET, PUT, DELETE for /api/movies/id
	mux.HandleFunc(apiMovies+"/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.userIdentity(h.handleGetMovieById)(w, r)
			return
		} else if r.Method == http.MethodDelete {
			h.userIdentity(h.handleDeleteMovie)(w, r)
			return
		} else if r.Method == http.MethodPut {
			h.userIdentity(h.handleUpdateMovie)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	//GET for /api/movies/searchbytitle
	mux.HandleFunc(apiMovies+"/searchbytitle", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.userIdentity(h.handleSearchMoviesByTitle)(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	//GET for /api/movies/searchbyactor
	mux.HandleFunc(apiMovies+"/searchbyactor", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.userIdentity(h.handleSearchMoviesByActorName)(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
