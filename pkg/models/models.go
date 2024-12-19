package models

import (
	"errors"
	"fmt"
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

func (t *Todo) FormatCreated() string { //Форматирование даты Created в формат ММ ДД ГГГ
	year, month, day := t.Created.Date()
	return fmt.Sprintf("%v %d, %d", month, day, year)
}

func (t *Todo) FormatExpires() string { //Форматирование даты Expires в формат ММ ДД ГГГ
	year, month, day := t.Expires.Date()
	return fmt.Sprintf("%v %d, %d", month, day, year)
}
