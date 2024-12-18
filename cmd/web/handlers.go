package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/speeddem0n/todobox/pkg/models"
)

func (app *application) home(wr http.ResponseWriter, resp *http.Request) {

	if resp.URL.Path != "/" { // Проверяется, если текущий путь URL запроса точно совпадает с шаблоном "/".
		app.notFound(wr) // хелпер notFound для возвращения клиенту ошибки 404.
		return
	}

	todoSlice, err := app.todos.Latest()
	if err != nil {
		app.serverError(wr, err)
		return
	}

	app.render(wr, "home.page.tmpl", &templateData{ // Используем помощника render() для отображения шаблона.
		Todos: todoSlice,
	})
}

func (app *application) showTodo(wr http.ResponseWriter, resp *http.Request) {

	id, err := strconv.Atoi(resp.URL.Query().Get("id")) //Извлекаем значение параметра id из URL

	if err != nil || id < 1 {
		app.notFound(wr) // страница не найдена 404.
		return
	}

	todo, err := app.todos.Get(id) // Метод Get() из pSQL модели todos
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) { // Если подходящей записи не найдено - ошибка 404
			app.notFound(wr)
		} else {
			app.serverError(wr, err)
		}
		return
	}
	app.render(wr, "show.page.tmpl", &templateData{
		Todo: todo,
	})
}

func (app *application) newTodo(wr http.ResponseWriter, resp *http.Request) {
	if resp.Method == "POST" {

		err := resp.ParseForm()
		if err != nil {
			app.serverError(wr, err) // Используем помошник serverError для обработки ошибки
		}
		title := resp.FormValue("title")
		content := resp.FormValue("content")
		expires := resp.FormValue("expires")

		id, err := app.todos.Insert(title, content, expires) // Добавляем данные в БД
		if err != nil {
			app.serverError(wr, err)
			return
		}

		http.Redirect(wr, resp, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther) // Перенаправляем пользователя на созданную заметку
	} else {
		app.render(wr, "create.page.tmpl", nil)
	}
}

func (app *application) deleteTodo(wr http.ResponseWriter, resp *http.Request) { // Handler для удаления заметок
	id, err := strconv.Atoi(resp.URL.Query().Get("id"))
	if err != nil {
		app.notFound(wr)
		return
	}
	err = app.todos.Delete(id)
	if err != nil {
		app.serverError(wr, err)
		return
	}
	http.Redirect(wr, resp, "/", http.StatusSeeOther)

}
