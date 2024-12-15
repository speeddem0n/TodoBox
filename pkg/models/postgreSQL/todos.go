package pSQL

import (
	"database/sql"

	"github.com/speeddem0n/todobox/pkg/models"
)

type TodoModel struct {
	DB *sql.DB
}

func (m *TodoModel) Insert(title, content, expires string) (int, error) {
	return 0, nil
}
func (m *TodoModel) Get(id int) (*models.Todo, error) {
	return nil, nil
}

func (m *TodoModel) Latest() ([]*models.Todo, error) {
	return nil, nil
}
