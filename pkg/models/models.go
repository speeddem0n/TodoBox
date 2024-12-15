package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: Подходящей записи не найдено")

type Todo struct { // Типы данных верхнего уровня
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
