package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/wintersakura/todo-api/pkg/model"
	"strings"
)

type TodoListSQL struct {
	db *sqlx.DB
}

func NewTodoListSQL(db *sqlx.DB) *TodoListSQL {
	return &TodoListSQL{db: db}
}

func (r *TodoListSQL) Create(userId int, list model.TodoList) (int, error) {
	t, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	createListQuery := fmt.Sprintf("INSERT INTO %s (`title`, `description`) VALUES (?, ?)", todoListsTable)
	res, err := t.Exec(createListQuery, list.Title, list.Description)
	if err != nil {
		t.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		t.Rollback()
		return 0, err
	}

	createUsersListsQuery := fmt.Sprintf("INSERT INTO %s (`user_id`, `list_id`) VALUES (?, ?)", usersListsTable)
	_, err = t.Exec(createUsersListsQuery, userId, id)
	if err != nil {
		t.Rollback()
		return 0, err
	}

	return int(id), t.Commit()
}

func (r *TodoListSQL) GetAll(userId int) ([]model.TodoList, error) {
	var lists []model.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul ON tl.id = ul.list_id WHERE ul.user_id = ?",
		todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *TodoListSQL) GetById(userId, listId int) (model.TodoList, error) {
	var list model.TodoList

	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description 
									FROM %s tl 
									INNER JOIN %s ul ON tl.id = ul.list_id 
									WHERE ul.user_id = ? AND ul.list_id = ?`,
		todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}

func (r *TodoListSQL) Delete(userId, listId int) error {
	query := fmt.Sprintf("DELETE tl FROM %s tl JOIN %s ul ON tl.id = ul.list_id WHERE ul.user_id = ? AND ul.list_id = ?",
		todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *TodoListSQL) Update(userId, listId int, input model.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("`title`=?"))
		args = append(args, *input.Title)
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("`description`=?"))
		args = append(args, *input.Description)
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl JOIN %s ul ON tl.id = ul.list_id SET %s WHERE ul.list_id = ? AND ul.user_id = ?",
		todoListsTable, usersListsTable, setQuery)
	args = append(args, listId, userId)

	_, err := r.db.Exec(query, args...)
	return err
}
