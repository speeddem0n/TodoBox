package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(wr http.ResponseWriter, resp *http.Request) {

	if resp.URL.Path != "/" { // Проверяется, если текущий путь URL запроса точно совпадает с шаблоном "/".
		app.notFound(wr) // хелпер notFound для возвращения клиенту ошибки 404.
		return
	}

	files := []string{
		".\\ui\\html\\home.page.tmpl",
		".\\ui\\html\\base.layout.tmpl",
		".\\ui\\html\\footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(wr, err) // Использование helper serverError
		return
	}

	err = ts.Execute(wr, nil)
	if err != nil {
		app.serverError(wr, err) // Использование helper serverError
	}

}

func (app *application) showSnippet(wr http.ResponseWriter, resp *http.Request) {

	id, err := strconv.Atoi(resp.URL.Query().Get("id")) //Извлекаем значение параметра id из URL

	if err != nil || id < 1 {
		app.notFound(wr) // хелпер notFound для возвращения клиенту ошибки 404.
	}
	fmt.Fprintf(wr, "Отображение выбранной заметки с ID %d...", id) // Вставкв значения из id в строку ответа
}

func (app *application) createSnippet(wr http.ResponseWriter, resp *http.Request) {
	if resp.Method != http.MethodPost {
		wr.Header().Set("Allow", http.MethodPost)
		app.clientError(wr, http.StatusMethodNotAllowed) // Испольхуем помошник clientError для обработки ошибки
		return
	}
	title := "Test, todo" // Тестовые данные для ввода в заметку
	content := "I need to work more AAAAAAAAAAAA"
	expires := "7"

	id, err := app.todos.Insert(title, content, expires) // Добавляем данные в БД
	if err != nil {
		app.serverError(wr, err)
		return
	}

	http.Redirect(wr, resp, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther) // Перенаправляем пользователя на созданную заметку
}
