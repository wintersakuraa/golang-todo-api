package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/wintersakura/todo-api/pkg/model"
	"strings"
)

type TodoItemSQL struct {
	db *sqlx.DB
}

func NewTodoItemSQL(db *sqlx.DB) *TodoItemSQL {
	return &TodoItemSQL{db: db}
}

func (r *TodoItemSQL) Create(listId int, item model.TodoItem) (int, error) {
	t, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	createItemQuery := fmt.Sprintf("INSERT INTO %s (`title`, `description`) VALUES (?, ?)", todoItemsTable)
	res, err := t.Exec(createItemQuery, item.Title, item.Description)
	if err != nil {
		t.Rollback()
		return 0, err
	}

	itemId, err := res.LastInsertId()
	if err != nil {
		t.Rollback()
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (`list_id`, `item_id`) VALUES (?, ?)", listsItemsTable)
	_, err = t.Exec(createListItemsQuery, listId, itemId)
	if err != nil {
		t.Rollback()
		return 0, err
	}

	return int(itemId), t.Commit()
}

func (r *TodoItemSQL) GetAll(userId, listId int) ([]model.TodoItem, error) {
	var items []model.TodoItem

	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti 
         							INNER JOIN %s li ON li.item_id = ti.id 
            					 	INNER JOIN %s ul ON ul.list_id = li.list_id
            					 	WHERE li.list_id = ? AND ul.user_id = ?`,
		todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errors.New("no items in the list")
	}

	return items, nil
}

func (r *TodoItemSQL) GetById(userId, itemId int) (model.TodoItem, error) {
	var item model.TodoItem

	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti 
         							INNER JOIN %s li ON li.item_id = ti.id 
            					 	INNER JOIN %s ul ON ul.list_id = li.list_id
            					 	WHERE ti.id = ? AND ul.user_id = ?`,
		todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Get(&item, query, itemId, userId); err != nil {
		return item, err
	}

	return item, nil
}

func (r *TodoItemSQL) Delete(userId, itemId int) error {
	query := fmt.Sprintf(`DELETE ti FROM %s ti 
									JOIN %s li ON ti.id = li.item_id 
									JOIN %s ul ON li.list_id = ul.list_id 
									WHERE ul.user_id = ? AND ti.id = ?`,
		todoItemsTable, listsItemsTable, usersListsTable)

	_, err := r.db.Exec(query, userId, itemId)
	return err
}

func (r *TodoItemSQL) Update(userId, itemId int, input model.UpdateItemInput) error {
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

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("`done`=?"))
		args = append(args, *input.Done)
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s ti 
    								JOIN %s li ON ti.id = li.item_id
    								JOIN %s ul ON li.list_id = ul.list_id 
    								SET %s WHERE ul.user_id = ? AND ti.id = ?`,
		todoItemsTable, listsItemsTable, usersListsTable, setQuery)
	args = append(args, userId, itemId)

	_, err := r.db.Exec(query, args...)
	return err
}
