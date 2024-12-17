package pSQL

import (
	"database/sql"
	"errors"

	"github.com/speeddem0n/todobox/pkg/models"
)

type TodoModel struct {
	DB *sql.DB
}

func (m *TodoModel) Insert(title, content, expires string) (int, error) { // Метод добавления новой записи в БД
	stmt := `INSERT INTO todo (title, content, created, expires)
	VALUES($1,$2, date_trunc('second', current_timestamp), 
	date_trunc('second', current_timestamp) + INTERVAL '1 day' * $3) RETURNING id` // SQL запрос для добавления новой заметки, возвращает Id последнего запроса

	var id int // Переменная для последующей записи id
	err := m.DB.QueryRow(stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return int(id), nil // ID имеет тип Int64, поэтому он конвертируется в int
}

func (m *TodoModel) Get(id int) (*models.Todo, error) {
	stmt := `SELECT id, title, content, created, expires FROM todo
	WHERE expires > current_timestamp AND id = $1` // SQL запрос для получения записи по ее ID

	row := m.DB.QueryRow(stmt, id) // QueryRow() для выолнения SQL запроса возвращает указатель  на объект sql.Row который содержит данные записи

	s := &models.Todo{} // Инициализируем указатель на новую структуру Todo

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires) // row.Scan Для копирования значения каждого поля в структуру Todo
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord // Если запрос был выполнен с ошибкой и ошибка обнаружена возвращаем ошибку из содели models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil // Возвращаем объект Todo

}

func (m *TodoModel) Latest() ([]*models.Todo, error) {
	return nil, nil
}
