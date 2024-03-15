package filmoteka

import (
	"errors"
)

type Actors struct {
	Id          int    `json:"id" db:"id"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	Gender      string `json:"gender" db:"gender"`
	DateOfBirth string `json:"date_of_birth" db:"date_of_birth"`
}

type Movies struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	ReleaseDate string `json:"release_date" db:"release_date"`
	Rating      int    `json:"rating" db:"rating"`
}

type ActorsWithMovies struct {
	Id          int    `json:"id" db:"id"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	Gender      string `json:"gender" db:"gender"`
	DateOfBirth string `json:"date_of_birth" db:"date_of_birth"`
	Movies      string `json:"movies" db:"movies"`
}

type MoviesWithActors struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	ReleaseDate string `json:"release_date" db:"release_date"`
	Rating      int    `json:"rating" db:"rating"`
	Actors      string `json:"actors" db:"actors"`
}

type UpdateActors struct {
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	Gender      *string `json:"gender"`
	DateOfBirth *string `json:"date_of_birth"`
}

type UpdateMovies struct {
	Id          *int    `json:"id" db:"id"`
	Title       *string `json:"title" db:"title"`
	Description *string `json:"description" db:"description"`
	ReleaseDate *string `json:"release_date" db:"release_date"`
	Rating      *int    `json:"rating" db:"rating"`
	Actors      *[]int  `json:"actors" db:"actors"`
}

func (u UpdateActors) Validate() error {
	if u.FirstName == nil && u.LastName == nil && u.Gender == nil && u.DateOfBirth == nil {
		return errors.New("update structure has no values")
	}
	return nil
}

func (u UpdateMovies) Validate() error {
	if u.Title == nil && u.Description == nil && u.ReleaseDate == nil && u.Rating == nil && u.Actors == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
