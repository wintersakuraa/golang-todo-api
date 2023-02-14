package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/wintersakura/todo-api/pkg/model"
)

type AuthSQL struct {
	db *sqlx.DB
}

func NewAuthSQL(db *sqlx.DB) *AuthSQL {
	return &AuthSQL{db: db}
}

func (r *AuthSQL) CreateUser(user model.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (`name`, username, password_hash) VALUES (?, ?, ?);", usersTable)

	res, err := r.db.Exec(query, user.Name, user.Username, user.Password)
	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()

	return int(id), nil
}

func (r *AuthSQL) GetUser(username, password string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT `id` FROM %s WHERE `username`=? AND `password_hash`=?", usersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}
