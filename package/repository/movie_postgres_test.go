package repository

import (
	"regexp"
	"testing"
	filmoteka "vk_restAPI"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestMoviePostgres_GetMovies(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' wasn't expected when opening a stub db connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewMoviePostgres(sqlxDB)

	expectedMovies := []filmoteka.MoviesWithActors{
		{
			Id:          1,
			Title:       "testTitle",
			Description: "testDescription",
			ReleaseDate: "1996-06-20",
			Rating:      7,
			Actors:      "Actor 1, Actor 2",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "description", "release_date", "rating", "actors"}).
		AddRow(expectedMovies[0].Id, expectedMovies[0].Title, expectedMovies[0].Description, expectedMovies[0].ReleaseDate, expectedMovies[0].Rating, expectedMovies[0].Actors)
	mock.ExpectQuery(regexp.QuoteMeta(`
    SELECT 
        m.id, 
        m.title, 
        m.description, 
        TO_CHAR (m.release_date, 'YYYY-MM-DD') AS release_date, 
        m.rating, 
        array_agg(a.first_name || ' ' || a.last_name) AS actors
    FROM 
        movies m
    LEFT JOIN 
        moviesactors ma ON m.id = ma.movie_id
    LEFT JOIN 
        actors a ON ma.actor_id = a.id
    GROUP BY 
        m.id, m.title, m.description, m.release_date, m.rating
    ORDER BY 
        m.rating DESC
    `)).WillReturnRows(rows)

	movies, err := repo.GetMovies()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	assert.NoError(t, err)
	assert.Equal(t, len(expectedMovies), len(movies))
	assert.Equal(t, expectedMovies[0].Title, movies[0].Title)
	assert.Equal(t, expectedMovies[0].Description, movies[0].Description)
	assert.Equal(t, expectedMovies[0].ReleaseDate, movies[0].ReleaseDate)
	assert.Equal(t, expectedMovies[0].Rating, movies[0].Rating)
	assert.Equal(t, expectedMovies[0].Actors, movies[0].Actors)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMoviePostgres_GetMoviesSortedByTitle(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' wasn't expected when opening a stub db connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewMoviePostgres(sqlxDB)

	// Ожидаемые результаты
	expectedMovies := []filmoteka.MoviesWithActors{
		{
			Id:          1,
			Title:       "First Title",
			Description: "First Description",
			ReleaseDate: "1996-06-20",
			Rating:      7,
			Actors:      "Actor 1, Actor 2",
		},
		{
			Id:          2,
			Title:       "Second Title",
			Description: "Second Description",
			ReleaseDate: "2000-01-01",
			Rating:      8,
			Actors:      "Actor 3, Actor 4",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "description", "release_date", "rating", "actors"}).
		AddRow(expectedMovies[0].Id, expectedMovies[0].Title, expectedMovies[0].Description, expectedMovies[0].ReleaseDate, expectedMovies[0].Rating, expectedMovies[0].Actors).
		AddRow(expectedMovies[1].Id, expectedMovies[1].Title, expectedMovies[1].Description, expectedMovies[1].ReleaseDate, expectedMovies[1].Rating, expectedMovies[1].Actors)
	mock.ExpectQuery(regexp.QuoteMeta(`
    SELECT 
        m.id, 
        m.title, 
        m.description, 
        TO_CHAR (m.release_date, 'YYYY-MM-DD') AS release_date, 
        m.rating, 
        array_agg(a.first_name || ' ' || a.last_name) AS actors
    FROM 
        movies m
    LEFT JOIN 
        moviesactors ma ON m.id = ma.movie_id
    LEFT JOIN 
        actors a ON ma.actor_id = a.id
    GROUP BY 
        m.id, m.title, m.description, m.release_date, m.rating
    ORDER BY 
        m.title ASC
    `)).WillReturnRows(rows)

	movies, err := repo.GetMoviesSortedByTitle()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	assert.NoError(t, err)
	assert.Equal(t, len(expectedMovies), len(movies))

	for i := range expectedMovies {
		assert.Equal(t, expectedMovies[i].Title, movies[i].Title)
		assert.Equal(t, expectedMovies[i].Description, movies[i].Description)
		assert.Equal(t, expectedMovies[i].ReleaseDate, movies[i].ReleaseDate)
		assert.Equal(t, expectedMovies[i].Rating, movies[i].Rating)
		assert.Equal(t, expectedMovies[i].Actors, movies[i].Actors)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMoviePostgres_GetMoviesSortedByDate(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' wasn't expected when opening a stub db connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewMoviePostgres(sqlxDB)

	expectedMovies := []filmoteka.MoviesWithActors{
		{
			Id:          1,
			Title:       "First Title",
			Description: "First Description",
			ReleaseDate: "1996-06-20",
			Rating:      7,
			Actors:      "Actor 1, Actor 2",
		},
		{
			Id:          2,
			Title:       "Second Title",
			Description: "Second Description",
			ReleaseDate: "2000-01-01",
			Rating:      8,
			Actors:      "Actor 3, Actor 4",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "description", "release_date", "rating", "actors"}).
		AddRow(expectedMovies[0].Id, expectedMovies[0].Title, expectedMovies[0].Description, expectedMovies[0].ReleaseDate, expectedMovies[0].Rating, expectedMovies[0].Actors).
		AddRow(expectedMovies[1].Id, expectedMovies[1].Title, expectedMovies[1].Description, expectedMovies[1].ReleaseDate, expectedMovies[1].Rating, expectedMovies[1].Actors)
	mock.ExpectQuery(regexp.QuoteMeta(`
    SELECT 
        m.id, 
        m.title, 
        m.description, 
        TO_CHAR (m.release_date, 'YYYY-MM-DD') AS release_date, 
        m.rating, 
        array_agg(a.first_name || ' ' || a.last_name) AS actors
    FROM 
        movies m
    LEFT JOIN 
        moviesactors ma ON m.id = ma.movie_id
    LEFT JOIN 
        actors a ON ma.actor_id = a.id
    GROUP BY 
        m.id, m.title, m.description, m.release_date, m.rating
    ORDER BY 
        m.release_date ASC
    `)).WillReturnRows(rows)

	movies, err := repo.GetMoviesSortedByDate()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	assert.NoError(t, err)
	assert.Equal(t, len(expectedMovies), len(movies))

	for i := range expectedMovies {
		assert.Equal(t, expectedMovies[i].Title, movies[i].Title)
		assert.Equal(t, expectedMovies[i].Description, movies[i].Description)
		assert.Equal(t, expectedMovies[i].ReleaseDate, movies[i].ReleaseDate)
		assert.Equal(t, expectedMovies[i].Rating, movies[i].Rating)
		assert.Equal(t, expectedMovies[i].Actors, movies[i].Actors)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMoviePostgres_GetMovieById(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' wasn't expected when opening a stub db connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewMoviePostgres(sqlxDB)

	expectedMovie := filmoteka.MoviesWithActors{
		Id:          1,
		Title:       "Test Title",
		Description: "Test Description",
		ReleaseDate: "1996-06-20",
		Rating:      7,
		Actors:      "Actor 1, Actor 2",
	}

	rows := sqlmock.NewRows([]string{"id", "title", "description", "release_date", "rating", "actors"}).
		AddRow(expectedMovie.Id, expectedMovie.Title, expectedMovie.Description, expectedMovie.ReleaseDate, expectedMovie.Rating, expectedMovie.Actors)
	mock.ExpectQuery(regexp.QuoteMeta(`
    SELECT 
        m.id, 
        m.title, 
        m.description, 
        TO_CHAR (m.release_date, 'YYYY-MM-DD') AS release_date, 
        m.rating, 
        array_agg(a.first_name || ' ' || a.last_name) AS actors
    FROM 
        movies m
    LEFT JOIN 
        moviesactors ma ON m.id = ma.movie_id
    LEFT JOIN 
        actors a ON ma.actor_id = a.id
    WHERE 
        m.id=$1
    GROUP BY 
        m.id, m.title, m.description, m.release_date, m.rating
    `)).WithArgs(expectedMovie.Id).WillReturnRows(rows)

	movie, err := repo.GetMovieById(expectedMovie.Id)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedMovie.Id, movie.Id)
	assert.Equal(t, expectedMovie.Title, movie.Title)
	assert.Equal(t, expectedMovie.Description, movie.Description)
	assert.Equal(t, expectedMovie.ReleaseDate, movie.ReleaseDate)
	assert.Equal(t, expectedMovie.Rating, movie.Rating)
	assert.Equal(t, expectedMovie.Actors, movie.Actors)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMoviePostgres_SearchMoviesByTitle(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' wasn't expected when opening a stub db connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewMoviePostgres(sqlxDB)

	expectedMovies := []filmoteka.MoviesWithActors{
		{
			Id:          1,
			Title:       "Test Title 1",
			Description: "Test Description 1",
			ReleaseDate: "1996-06-20",
			Rating:      7,
			Actors:      "Actor 1, Actor 2",
		},
		{
			Id:          2,
			Title:       "Test Title 2",
			Description: "Test Description 2",
			ReleaseDate: "2000-01-01",
			Rating:      8,
			Actors:      "Actor 3, Actor 4",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "description", "release_date", "rating", "actors"}).
		AddRow(expectedMovies[0].Id, expectedMovies[0].Title, expectedMovies[0].Description, expectedMovies[0].ReleaseDate, expectedMovies[0].Rating, expectedMovies[0].Actors).
		AddRow(expectedMovies[1].Id, expectedMovies[1].Title, expectedMovies[1].Description, expectedMovies[1].ReleaseDate, expectedMovies[1].Rating, expectedMovies[1].Actors)
	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT 
            m.*,
            array_agg(concat(a.first_name, ' ', a.last_name)) AS actors
        FROM 
            movies m
        INNER JOIN 
            moviesactors ma ON m.id = ma.movie_id
        INNER JOIN 
            actors a ON ma.actor_id = a.id
        WHERE 
            LOWER(m.title) LIKE LOWER($1)
        GROUP BY 
            m.id
    `)).WithArgs("%test%").WillReturnRows(rows)

	movies, err := repo.SearchMoviesByTitle("test")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	assert.NoError(t, err)
	assert.Equal(t, len(expectedMovies), len(movies))
	for i := range expectedMovies {
		assert.Equal(t, expectedMovies[i].Id, movies[i].Id)
		assert.Equal(t, expectedMovies[i].Title, movies[i].Title)
		assert.Equal(t, expectedMovies[i].Description, movies[i].Description)
		assert.Equal(t, expectedMovies[i].ReleaseDate, movies[i].ReleaseDate)
		assert.Equal(t, expectedMovies[i].Rating, movies[i].Rating)
		assert.Equal(t, expectedMovies[i].Actors, movies[i].Actors)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMoviePostgres_SearchMovieByActorName(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' wasn't expected when opening a stub db connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewMoviePostgres(sqlxDB)

	expectedMovie := filmoteka.MoviesWithActors{
		Id:          1,
		Title:       "Test Title 1",
		Description: "Test Description 1",
		ReleaseDate: "1996-06-20",
		Rating:      7,
		Actors:      "Actor 1, Actor 2",
	}

	rows := sqlmock.NewRows([]string{"id", "title", "description", "release_date", "rating", "actors"}).
		AddRow(expectedMovie.Id, expectedMovie.Title, expectedMovie.Description, expectedMovie.ReleaseDate, expectedMovie.Rating, expectedMovie.Actors)
	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT 
            m.*,
            array_agg(concat(a.first_name, ' ', a.last_name)) AS actors
        FROM 
            movies m
        INNER JOIN 
            moviesactors ma ON m.id = ma.movie_id
        INNER JOIN 
            actors a ON ma.actor_id = a.id
        WHERE 
            LOWER(a.first_name) LIKE LOWER($1) OR LOWER(a.last_name) LIKE LOWER($1)
        GROUP BY 
            m.id
    `)).WithArgs("%actor%").WillReturnRows(rows)

	movies, err := repo.SearchMovieByActorName("actor")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	assert.NoError(t, err)
	assert.Equal(t, 1, len(movies))
	assert.Equal(t, expectedMovie.Id, movies[0].Id)
	assert.Equal(t, expectedMovie.Title, movies[0].Title)
	assert.Equal(t, expectedMovie.Description, movies[0].Description)
	assert.Equal(t, expectedMovie.ReleaseDate, movies[0].ReleaseDate)
	assert.Equal(t, expectedMovie.Rating, movies[0].Rating)
	assert.Equal(t, expectedMovie.Actors, movies[0].Actors)

	assert.NoError(t, mock.ExpectationsWereMet())
}
