package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(wr http.ResponseWriter, err error) { //Оборачиваем ошибки сервера
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(wr http.ResponseWriter, status int) { //Оборачиваем клиентские ошибки
	http.Error(wr, http.StatusText(status), status)
}

func (app *application) notFound(wr http.ResponseWriter) { // Обертка для ошибки 404 not found
	app.clientError(wr, http.StatusNotFound)
}

func (app *application) render(wr http.ResponseWriter, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	/* Извлекаем соответствующий набор шаблонов из кэша в зависимости от названия страницы
	   (например, 'home.page.tmpl'). Если в кэше нет записи запрашиваемого шаблона, то
	   вызывается вспомогательный метод serverError()*/
	if !ok {
		app.serverError(wr, fmt.Errorf("Шаблон %s не существует!", name))
		return
	}

	err := ts.Execute(wr, td) // Рендерим файлы шаблона, передавая динамические данные из переменной `td`.
	if err != nil {
		app.serverError(wr, err)
	}
}
