package service

import (
	filmoteka "vk_restAPI"
	"vk_restAPI/package/repository"
)

type Authorization interface {
	CreateUser(user filmoteka.User) (int, error)
	GetUserStatus(id int) (bool, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Actors interface {
	CreateActor(actor filmoteka.Actors) (int, error)
	DeleteActor(actorId int) error
	UpdateActor(actorId int, input filmoteka.UpdateActors) error
}

type Movies interface {
	CreateMovie(userId int, movie filmoteka.Movies, actorIDs []int) (int, error)
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

type Service struct {
	Authorization
	Actors
	Movies
	MoviesWithActors
	ActorsWithMovies
}

// Service access databaseses
func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization:    NewAuthService(repos.Authorization),
		Actors:           NewActorService(repos.Actors),
		Movies:           NewMovieService(repos.Movies),
		MoviesWithActors: NewMoviesWithActorsService(repos.MoviesWithActors),
		ActorsWithMovies: NewActorsWithMoviesService(repos.ActorsWithMovies),
	}
}
