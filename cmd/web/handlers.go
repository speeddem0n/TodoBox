package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/speeddem0n/todobox/pkg/models"
)

// Здесь находятся Handler'ы всех путей web приложения

func (app *application) home(wr http.ResponseWriter, resp *http.Request) { // Handler для главной страницы "/"

	if resp.URL.Path != "/" { // Проверяется, если текущий путь URL запроса точно совпадает с шаблоном "/".
		app.notFound(wr) // Хелпер notFound для возвращения клиенту ошибки 404.
		return
	}

	todoSlice, err := app.todos.Latest() // Методом Latest из модели БД достаем последние 10 заметок и записываем их в переменную
	if err != nil {
		app.serverError(wr, err) // Используем помошник serverError для обработки ошибки
		return
	}

	app.render(wr, "home.page.tmpl", &templateData{ // Используем помощника render() для отображения шаблона.
		Todos: todoSlice,
	})
}

func (app *application) showTodo(wr http.ResponseWriter, resp *http.Request) { // Handler для отображения заметки по ее ID "/todo"

	id, err := strconv.Atoi(resp.URL.Query().Get("id")) //Извлекаем значение параметра id из URL

	if err != nil || id < 1 {
		app.notFound(wr) // Страница не найдена 404.
		return
	}

	todo, err := app.todos.Get(id) // Метод Get() из pSQL модели todos
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) { // Если подходящей записи не найдено - ошибка 404
			app.notFound(wr)
		} else {
			app.serverError(wr, err) // Используем помошник serverError для обработки ошибки
		}
		return
	}
	app.render(wr, "show.page.tmpl", &templateData{ // Отоброшаем страницу с заметкой
		Todo: todo,
	})
}

func (app *application) newTodo(wr http.ResponseWriter, resp *http.Request) { // Handler для создания новой заметки "/create"
	if resp.Method == "POST" { // Если метод POST...

		err := resp.ParseForm() // Парсим значения из URL
		if err != nil {
			app.serverError(wr, err) // Используем помошник serverError для обработки ошибки
		}
		title := resp.FormValue("title")
		content := resp.FormValue("content")
		expires := resp.FormValue("expires")

		id, err := app.todos.Insert(title, content, expires) // Добавляем данные в БД
		if err != nil {
			app.serverError(wr, err) // Обрабатываем ошибку со стороны сервера если она есть
			return
		}

		http.Redirect(wr, resp, fmt.Sprintf("/todo?id=%d", id), http.StatusSeeOther) // Перенаправляем пользователя на созданную заметку
	} else {
		app.render(wr, "create.page.tmpl", nil) // Отображается страница для создания заметки
	}
}

func (app *application) deleteTodo(wr http.ResponseWriter, resp *http.Request) { // Handler для удаления заметок "/delete"
	id, err := strconv.Atoi(resp.URL.Query().Get("id")) // Достаем query id из URL и преобразоввываем его в тип INT
	if err != nil {
		app.notFound(wr) // Ошибка 404 если такого id не существует
		return
	}
	err = app.todos.Delete(id) // Используем метод Delete() из модели БД для удаления записи по ее id
	if err != nil {
		app.serverError(wr, err) // Обрабатываем ошибку со стороны сервера если она есть
		return
	}
	http.Redirect(wr, resp, "/", http.StatusSeeOther) // Перенаправляем юзера обратно на главную страницу

}

func (app *application) searchTodo(wr http.ResponseWriter, resp *http.Request) { // Handler для поиска заметок по ID /search
	err := resp.ParseForm() // Парсим значения из URL
	if err != nil {
		app.serverError(wr, err) // Используем помошник serverError для обработки ошибки
		return
	}
	searchID := resp.FormValue("searchID") // Получаем ID от пользователя
	id, err := strconv.Atoi(searchID)      // Конвертируем ID в INT
	if err != nil {
		app.notFound(wr) // Ошибка 404 если такого ID не существует
		return
	}

	http.Redirect(wr, resp, fmt.Sprintf("/todo?id=%d", id), http.StatusSeeOther) // Перенаправляем пользователя на страницу с заметкой если все прошло успешно
}
