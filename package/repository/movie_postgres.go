package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	filmoteka "vk_restAPI"

	"github.com/jmoiron/sqlx"
)

type MoviePostgres struct {
	db *sqlx.DB
}

func NewMoviePostgres(db *sqlx.DB) *MoviePostgres {
	return &MoviePostgres{db: db}
}

func (m *MoviePostgres) CreateMovie(movie filmoteka.Movies, actorIDs []int) (int, error) {
	var id int

	var exisitngID int

	query := fmt.Sprintf("SELECT id FROM %s WHERE title = $1 AND description = $2", moviesTable)
	row := m.db.QueryRow(query, movie.Title, movie.Description)
	if err := row.Scan(&exisitngID); err == nil {
		return exisitngID, errors.New("movie with the same parameters already exists")
	} else if err != sql.ErrNoRows {
		return 0, err
	}

	query = fmt.Sprintf("INSERT INTO %s (title, description, release_date, rating) VALUES ($1, $2, $3, $4) RETURNING id", moviesTable)
	row = m.db.QueryRow(query, movie.Title, movie.Description, movie.ReleaseDate, movie.Rating)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	for _, actroID := range actorIDs {
		query = fmt.Sprintf("INSERT INTO %s (movie_id, actor_id) VALUES ($1, $2)", moviesActorsTable)
		_, err := m.db.Exec(query, id, actroID)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (m *MoviePostgres) GetMovies() ([]filmoteka.MoviesWithActors, error) {
	var movies []filmoteka.MoviesWithActors
	query := fmt.Sprintf(`
    SELECT 
        m.id, 
        m.title, 
        m.description, 
        TO_CHAR (m.release_date, 'YYYY-MM-DD') AS release_date, 
        m.rating, 
        array_agg(a.first_name || ' ' || a.last_name) AS actors
    FROM 
        %s m
    LEFT JOIN 
        %s ma ON m.id = ma.movie_id
    LEFT JOIN 
        %s a ON ma.actor_id = a.id
    GROUP BY 
        m.id, m.title, m.description, m.release_date, m.rating
	ORDER BY 
		m.rating DESC
`, moviesTable, moviesActorsTable, actorsTable)
	err := m.db.Select(&movies, query)
	return movies, err

}

func (m *MoviePostgres) GetMoviesSortedByTitle() ([]filmoteka.MoviesWithActors, error) {
	var movies []filmoteka.MoviesWithActors
	query := fmt.Sprintf(`
    SELECT 
        m.id, 
        m.title, 
        m.description, 
        TO_CHAR (m.release_date, 'YYYY-MM-DD') AS release_date, 	
        m.rating, 
        array_agg(a.first_name || ' ' || a.last_name) AS actors
    FROM 
        %s m
    LEFT JOIN 
        %s ma ON m.id = ma.movie_id
    LEFT JOIN 
        %s a ON ma.actor_id = a.id
    GROUP BY 
        m.id, m.title, m.description, m.release_date, m.rating
	ORDER BY 
		m.title ASC
`, moviesTable, moviesActorsTable, actorsTable)
	err := m.db.Select(&movies, query)
	return movies, err

}

func (m *MoviePostgres) GetMoviesSortedByDate() ([]filmoteka.MoviesWithActors, error) {
	var movies []filmoteka.MoviesWithActors
	query := fmt.Sprintf(`
    SELECT 
        m.id, 
        m.title, 
        m.description, 
        TO_CHAR (m.release_date, 'YYYY-MM-DD') AS release_date, 			 
        m.rating, 
        array_agg(a.first_name || ' ' || a.last_name) AS actors
    FROM 
        %s m
    LEFT JOIN 
        %s ma ON m.id = ma.movie_id
    LEFT JOIN 
        %s a ON ma.actor_id = a.id
    GROUP BY 
        m.id, m.title, m.description, m.release_date, m.rating
	ORDER BY 
		m.release_date ASC
`, moviesTable, moviesActorsTable, actorsTable)
	err := m.db.Select(&movies, query)
	return movies, err

}

func (m *MoviePostgres) GetMovieById(movieId int) (filmoteka.MoviesWithActors, error) {
	var movie filmoteka.MoviesWithActors
	query := fmt.Sprintf(`
    SELECT 
        m.id, 
        m.title, 
        m.description, 
        TO_CHAR (m.release_date, 'YYYY-MM-DD') AS release_date, 
        m.rating, 
        array_agg(a.first_name || ' ' || a.last_name) AS actors
    FROM 
        %s m
    LEFT JOIN 
        %s ma ON m.id = ma.movie_id
    LEFT JOIN 
        %s a ON ma.actor_id = a.id
	WHERE 
		m.id=$1
    GROUP BY 
        m.id, m.title, m.description, m.release_date, m.rating
`, moviesTable, moviesActorsTable, actorsTable)
	err := m.db.Get(&movie, query, movieId)
	return movie, err

}

func (m *MoviePostgres) DeleteMovie(movieId int) error {

	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := fmt.Sprintf("DELETE FROM %s WHERE movie_id=$1", moviesActorsTable)
	if _, err := m.db.Exec(query, movieId); err != nil {
		return err
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE id=$1", moviesTable)
	if _, err = m.db.Exec(query, movieId); err != nil {
		return err
	}

	return nil
}

func (m *MoviePostgres) UpdateMovie(movieId int, input filmoteka.UpdateMovies) error {
	setValue := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValue = append(setValue, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValue = append(setValue, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.ReleaseDate != nil {
		setValue = append(setValue, fmt.Sprintf("release_date=$%d", argId))
		args = append(args, *input.ReleaseDate)
		argId++
	}

	if input.Rating != nil {
		setValue = append(setValue, fmt.Sprintf("rating=$%d", argId))
		args = append(args, *input.Rating)
		argId++
	}

	if input.Actors != nil {

		deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE movie_id=$1", moviesActorsTable)
		_, err := m.db.Exec(deleteQuery, movieId)
		if err != nil {
			return err
		}

		for _, actorID := range *input.Actors {
			insertQuery := fmt.Sprintf("INSERT INTO %s (actor_id, movie_id) VALUES ($1, $2)", moviesActorsTable)
			_, err := m.db.Exec(insertQuery, actorID, movieId)
			if err != nil {
				return err
			}
		}
	}

	setQuery := strings.Join(setValue, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", moviesTable, setQuery, argId)
	args = append(args, movieId)

	log.Printf("updateQuerry: %s", query)
	log.Printf("args :%v", args)

	_, err := m.db.Exec(query, args...)
	return err

}

func (m *MoviePostgres) SearchMoviesByTitle(fragment string) ([]filmoteka.MoviesWithActors, error) {
	var movies []filmoteka.MoviesWithActors

	query := fmt.Sprintf(`
		SELECT 
			m.*,
			array_agg(concat(a.first_name, ' ', a.last_name)) AS actors
		FROM 
			%s m
		INNER JOIN 
			%s ma ON m.id = ma.movie_id
		INNER JOIN 
			%s a ON ma.actor_id = a.id
		WHERE 
			LOWER(m.title) LIKE LOWER($1)
		GROUP BY 
			m.id
	`, moviesTable, moviesActorsTable, actorsTable)

	fragment = "%" + fragment + "%"

	err := m.db.Select(&movies, query, fragment)
	if err != nil {
		return nil, errors.New("sql error")
	}

	return movies, nil
}

func (m *MoviePostgres) SearchMovieByActorName(fragment string) ([]filmoteka.MoviesWithActors, error) {
	var movies []filmoteka.MoviesWithActors

	query := fmt.Sprintf(`
		SELECT 
			m.*,
			array_agg(concat(a.first_name, ' ', a.last_name)) AS actors
		FROM 
			%s m
		INNER JOIN 
			%s ma ON m.id = ma.movie_id
		INNER JOIN 
			%s a ON ma.actor_id = a.id
		WHERE 
			LOWER(a.first_name) LIKE LOWER($1) OR LOWER(a.last_name) LIKE LOWER($1)
		GROUP BY 
			m.id
	`, moviesTable, moviesActorsTable, actorsTable)

	fragment = "%" + fragment + "%"

	err := m.db.Select(&movies, query, fragment)
	if err != nil {
		return nil, err
	}

	return movies, nil
}
