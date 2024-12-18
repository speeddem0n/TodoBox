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

func (m *TodoModel) Get(id int) (*models.Todo, error) { // Метод возвращает заметку по ее ID
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

func (m *TodoModel) Latest() ([]*models.Todo, error) { // Метод возвращает последние 10 заметок
	stmt := `SELECT id, title, content, created, expires FROM todo
	WHERE expires > current_timestamp
	ORDER BY created DESC
	LIMIT 10` // SQL запрос

	rows, err := m.DB.Query(stmt) // Query() для выполнения SQL запроса. в ответ получим sql.Rows с результатами запроса
	if err != nil {
		return nil, err
	}

	defer rows.Close() // Откладываем вызов rows.Close(), чтобы быть уверенным, что набор результатов из sql.Rows

	var todos []*models.Todo // Инициализация пустого среза для хранения объектов models.Todo

	for rows.Next() { // Next() для перебора результата
		s := &models.Todo{}                                                  // Указатель на новую структуру Todo
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires) // Scan(), чтобы скопировать значения полей в структуру.
		if err != nil {
			return nil, err
		}
		todos = append(todos, s) // Добовляем структуру в срез
	}

	if err = rows.Err(); err != nil { // Вызываем метод rows.Err(), чтобы узнать если в ходе работы у нас не возникла какая либо ошибка
		return nil, err
	}

	return todos, nil // Возвращаем срез с данными
}

func (m *TodoModel) Delete(id int) error {
	stmt := `DELETE FROM todo where id = $1`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}
