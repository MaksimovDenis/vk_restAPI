package repository

import (
	filmoteka "vk_restAPI"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user filmoteka.User) (int, error)
	GetUser(username, password string) (filmoteka.User, error)
	GetUserStatus(id int) (bool, error)
}

type Actors interface {
	CreateActor(actor filmoteka.Actors) (int, error)
	DeleteActor(actorId int) error
	UpdateActor(actorId int, input filmoteka.UpdateActors) error
}

type Movies interface {
	CreateMovie(movie filmoteka.Movies, actorIDs []int) (int, error)
	DeleteMovie(movieId int) error
}

type ActorsWithMovies interface {
	GetActors() ([]filmoteka.ActorsWithMovies, error)
	GetActorById(actorId int) (filmoteka.ActorsWithMovies, error)
}

type MoviesWithActors interface {
	GetMovies() ([]filmoteka.MoviesWithActors, error)
	GetMoviesSortedByTitle() ([]filmoteka.MoviesWithActors, error)
	GetMoviesSortedByDate() ([]filmoteka.MoviesWithActors, error)
	GetMovieById(movieId int) (filmoteka.MoviesWithActors, error)
	UpdateMovie(movieId int, input filmoteka.UpdateMovies) error
	SearchMoviesByTitle(fragment string) ([]filmoteka.MoviesWithActors, error)
	SearchMovieByActorName(fragment string) ([]filmoteka.MoviesWithActors, error)
}

type Repository struct {
	Authorization
	Actors
	Movies
	MoviesWithActors
	ActorsWithMovies
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization:    NewAuthPostgres(db),
		Actors:           NewActorPostgres(db),
		Movies:           NewMoviePostgres(db),
		MoviesWithActors: NewMoviePostgres(db),
		ActorsWithMovies: NewActorPostgres(db),
	}
}
