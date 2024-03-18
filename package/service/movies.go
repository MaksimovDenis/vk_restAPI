package service

import (
	filmoteka "vk_restAPI"
	"vk_restAPI/package/repository"
)

type MovieService struct {
	repo repository.Movies
}

type MoviesWithActorsService struct {
	repo repository.MoviesWithActors
}

func NewMovieService(repo repository.Movies) *MovieService {
	return &MovieService{repo: repo}
}

func NewMoviesWithActorsService(repo repository.MoviesWithActors) *MoviesWithActorsService {
	return &MoviesWithActorsService{repo: repo}
}

func (m *MovieService) CreateMovie(movie filmoteka.Movies, actorIDs []int) (int, error) {
	return m.repo.CreateMovie(movie, actorIDs)
}

func (m *MoviesWithActorsService) GetMovies() ([]filmoteka.MoviesWithActors, error) {
	return m.repo.GetMovies()
}

func (m *MoviesWithActorsService) GetMoviesSortedByTitle() ([]filmoteka.MoviesWithActors, error) {
	return m.repo.GetMoviesSortedByTitle()
}

func (m *MoviesWithActorsService) GetMoviesSortedByDate() ([]filmoteka.MoviesWithActors, error) {
	return m.repo.GetMoviesSortedByDate()
}

func (m *MoviesWithActorsService) GetMovieById(movieId int) (filmoteka.MoviesWithActors, error) {
	return m.repo.GetMovieById(movieId)
}

func (m *MovieService) DeleteMovie(movieId int) error {
	return m.repo.DeleteMovie(movieId)
}

func (m *MoviesWithActorsService) UpdateMovie(movieId int, input filmoteka.UpdateMovies) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return m.repo.UpdateMovie(movieId, input)
}

func (m *MoviesWithActorsService) SearchMoviesByTitle(fragment string) ([]filmoteka.MoviesWithActors, error) {
	return m.repo.SearchMoviesByTitle(fragment)
}

func (m *MoviesWithActorsService) SearchMovieByActorName(fragment string) ([]filmoteka.MoviesWithActors, error) {
	return m.repo.SearchMovieByActorName(fragment)
}
