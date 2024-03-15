package repository

import (
	"fmt"
	filmoteka "vk_restAPI"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateUser(user filmoteka.User) (int, error) {
	var id int
	qurey := fmt.Sprintf("INSERT INTO %s (username, password_hash, is_admin) values ($1, $2, $3) RETURNING id", userTable)

	row := a.db.QueryRow(qurey, user.Username, user.Password, user.Is_admin)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (a *AuthPostgres) GetUser(username, password string) (filmoteka.User, error) {
	var user filmoteka.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", userTable)
	err := a.db.Get(&user, query, username, password)

	return user, err
}

func (a *AuthPostgres) GetUserStatus(id int) (bool, error) {
	var isAdmin bool
	query := fmt.Sprintf("SELECT is_admin FROM %s WHERE id=$1", userTable)
	err := a.db.Get(&isAdmin, query, id)

	return isAdmin, err
}
