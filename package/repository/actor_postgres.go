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

type ActorPostgres struct {
	db *sqlx.DB
}

func NewActorPostgres(db *sqlx.DB) *ActorPostgres {
	return &ActorPostgres{db: db}
}

func (a *ActorPostgres) CreateActor(actor filmoteka.Actors) (int, error) {
	var id int

	var exisitngID int

	query := fmt.Sprintf("SELECT id FROM %s WHERE first_name=$1 AND last_name=$2", actorsTable)
	row := a.db.QueryRow(query, actor.FirstName, actor.LastName)
	if err := row.Scan(&exisitngID); err == nil {
		return exisitngID, errors.New("actor with the same name already exists")
	} else if err != sql.ErrNoRows {
		return 0, err
	}

	query = fmt.Sprintf("INSERT INTO %s (first_name, last_name, gender, date_of_birth) VALUES ($1, $2, $3, $4) RETURNING id", actorsTable)
	row = a.db.QueryRow(query, actor.FirstName, actor.LastName, actor.Gender, actor.DateOfBirth)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (a *ActorPostgres) GetActors() ([]filmoteka.ActorsWithMovies, error) {
	var actors []filmoteka.ActorsWithMovies

	query := fmt.Sprintf(`
		SELECT 
			a.id, 
			a.first_name, 
			a.last_name, 
			a.gender, 
			a.date_of_birth, 
			array_agg(m.title) AS movies
		FROM 
			%s a
		LEFT JOIN 
			%s ma ON a.id = ma.actor_id
		LEFT JOIN 
			%s m ON ma.movie_id = m.id
		GROUP BY 
			a.id
	`, actorsTable, moviesActorsTable, moviesTable)

	err := a.db.Select(&actors, query)
	if err != nil {
		return nil, err
	}

	return actors, nil

}

func (a *ActorPostgres) GetActorById(actorId int) (filmoteka.ActorsWithMovies, error) {
	var actor filmoteka.ActorsWithMovies

	query := fmt.Sprintf(`
		SELECT 
			a.id, 
			a.first_name, 
			a.last_name, 
			a.gender, 
			a.date_of_birth, 
			array_agg(m.title) AS movies
		FROM 
			%s a
		LEFT JOIN 
			%s ma ON a.id = ma.actor_id
		LEFT JOIN 
			%s m ON ma.movie_id = m.id
		WHERE 
			a.id=$1
		GROUP BY 
			a.id
	`, actorsTable, moviesActorsTable, moviesTable)

	err := a.db.Get(&actor, query, actorId)
	return actor, err

}

func (a *ActorPostgres) DeleteActor(actorId int) error {

	qurey := fmt.Sprintf("DELETE FROM %s WHERE id=$1", actorsTable)
	_, err := a.db.Exec(qurey, actorId)

	return err
}

func (a *ActorPostgres) UpdateActor(actorId int, input filmoteka.UpdateActors) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.FirstName != nil {
		setValues = append(setValues, fmt.Sprintf("first_name=$%d", argId))
		args = append(args, *input.FirstName)
		argId++
	}

	if input.LastName != nil {
		setValues = append(setValues, fmt.Sprintf("last_name=$%d", argId))
		args = append(args, *input.LastName)
		argId++
	}

	if input.Gender != nil {
		setValues = append(setValues, fmt.Sprintf("gender=$%d", argId))
		args = append(args, *input.Gender)
		argId++
	}

	if input.DateOfBirth != nil {
		setValues = append(setValues, fmt.Sprintf("date_of_birth=$%d", argId))
		args = append(args, *input.DateOfBirth)
		argId++
	}
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", actorsTable, setQuery, argId)
	args = append(args, actorId)

	log.Printf("updateQuerry: %s", query)
	log.Printf("args :%v", args)

	_, err := a.db.Exec(query, args...)
	return err

}
