package repository

import (
	"testing"
	filmoteka "vk_restAPI"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestActorPostgres_CreateActor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' wasn't expected when opening a stub db connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	actorRepo := NewActorPostgres(sqlxDB)

	type args struct {
		atctor filmoteka.Actors
	}
	type mockBehaivior func(args args)

	testTable := []struct {
		name          string
		mockBehaivior mockBehaivior
		args          args
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				atctor: filmoteka.Actors{
					Id:          1,
					FirstName:   "Denis",
					LastName:    "Maksimov",
					Gender:      "male",
					DateOfBirth: "1996-06-20",
				},
			},

			mockBehaivior: func(args args) {

				mock.ExpectQuery("SELECT id FROM actors WHERE first_name=\\$1 AND last_name=\\$2").
					WithArgs("Denis", "Maksimov").
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				rows := sqlmock.NewRows([]string{"id"}).AddRow(args.atctor.Id)

				mock.ExpectQuery("INSERT INTO actors").
					WithArgs(args.atctor.FirstName, args.atctor.LastName, args.atctor.Gender, args.atctor.DateOfBirth).WillReturnRows(rows)
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaivior(testCase.args)

			got, err := actorRepo.CreateActor(testCase.args.atctor)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.args.atctor.Id, got)
			}
		})
	}

}

func TestActorPostgres_GetActorById(t *testing.T) {

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := NewActorPostgres(sqlx.NewDb(mockDB, "sqlmock"))

	//Ecpected result
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "gender", "date_of_birth", "movies"}).
		AddRow(1, "Denis", "Maksimov", "Male", "1996-06-20", "Movie 1, Movie 2")

	mock.ExpectQuery("^SELECT (.+) FROM actors a (.+)$").WithArgs(1).WillReturnRows(rows)

	actor, err := repo.GetActorById(1)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	assert.NoError(t, err)

	assert.Equal(t, 1, actor.Id)
	assert.Equal(t, "Denis", actor.FirstName)
	assert.Equal(t, "Maksimov", actor.LastName)
	assert.Equal(t, "Male", actor.Gender)
	assert.Equal(t, "1996-06-20", actor.DateOfBirth)
	assert.Equal(t, "Movie 1, Movie 2", actor.Movies)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActorPostgres_GetActors(t *testing.T) {

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ar error '%s' wasn't expected when opening a stub db connection", err)
	}
	defer mockDB.Close()

	repo := NewActorPostgres(sqlx.NewDb(mockDB, "sqlmock"))

	//Ecpected result
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "gender", "date_of_birth", "movies"}).
		AddRow(1, "Denis", "Maksimov", "Male", "1996-06-20", "Movie 1, Movie 2")

	mock.ExpectQuery("^SELECT (.+) FROM actors a (.+)$").WillReturnRows(rows)

	actors, err := repo.GetActors()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	assert.NoError(t, err)

	assert.Equal(t, 1, len(actors))
	assert.Equal(t, "Denis", actors[0].FirstName)
	assert.Equal(t, "Maksimov", actors[0].LastName)
	assert.Equal(t, "Male", actors[0].Gender)
	assert.Equal(t, "1996-06-20", actors[0].DateOfBirth)
	assert.Equal(t, "Movie 1, Movie 2", actors[0].Movies)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActorPostgres_DeleteActor(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' wasn't expected when opening a stub db connection", err)
	}
	defer mockDB.Close()

	repo := NewActorPostgres(sqlx.NewDb(mockDB, "sqlmock"))

	actorID := 1
	query := "DELETE FROM actors WHERE id=\\$1"
	mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 1)).WillReturnError(nil)

	err = repo.DeleteActor(actorID)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}
