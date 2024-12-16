package pSQL

import (
	"database/sql"

	"github.com/speeddem0n/todobox/pkg/models"
)

type TodoModel struct {
	DB *sql.DB
}

func (m *TodoModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO todo (title, content, created, expires)
	VALUES($1,$2, date_trunc('second', current_timestamp), 
	date_trunc('second', current_timestamp) + INTERVAL '1 day' * $3) RETURNING id`

	var id int
	err := m.DB.QueryRow(stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return int(id), nil // ID имеет тип Int64, поэтому он конвертируется в int
}

func (m *TodoModel) Get(id int) (*models.Todo, error) {
	return nil, nil
}

func (m *TodoModel) Latest() ([]*models.Todo, error) {
	return nil, nil
}
