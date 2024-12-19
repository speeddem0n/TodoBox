package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: Подходящей записи не найдено") // Ошибка если нет токой записи в БД

type Todo struct { // Типы данных верхнего уровня
	ID      int       // ID заметки
	Title   string    // Название заметки
	Content string    // Содержимое заметки
	Created time.Time // Дата создания
	Expires time.Time // Дата удаления
}
