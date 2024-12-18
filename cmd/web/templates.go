package main

import (
	"github.com/speeddem0n/todobox/pkg/models"
)

type templateData struct { // Хранилище для любых динамических данных, которые нужно передать HTML-шаблонам.
	Todo *models.Todo
}
