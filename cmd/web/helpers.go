package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Здесь обернуты различные ошибки и метод render для отобрашения html шаблонов

func (app *application) serverError(wr http.ResponseWriter, err error) { // Оборачиваем ошибки сервера
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack()) // Выводим саму ошибку а потом весь стэк для дебага
	app.errorLog.Output(2, trace)

	http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // Отоброжение ошибки сервера в браузере
}

func (app *application) clientError(wr http.ResponseWriter, status int) { // Оборачиваем клиентские ошибки
	http.Error(wr, http.StatusText(status), status) // Отоброжение клиентской ошибки в браузере
}

func (app *application) notFound(wr http.ResponseWriter) { // Обертка для ошибки 404 not found
	app.clientError(wr, http.StatusNotFound) // Отоброжение ошибки 404 not found в браузере
}

func (app *application) render(wr http.ResponseWriter, name string, td *templateData) { // Помошник для отображения нужных html шаблонов для страницы
	ts, ok := app.templateCache[name]
	/* Извлекаем соответствующий набор шаблонов из кэша в зависимости от названия страницы
	   (например, 'home.page.tmpl'). Если в кэше нет записи запрашиваемого шаблона, то
	   вызывается вспомогательный метод serverError()*/
	if !ok {
		app.serverError(wr, fmt.Errorf("Шаблон %s не существует!", name)) // Ошибка сервера если такого шаблона не существует
		return
	}

	err := ts.Execute(wr, td) // Рендерим файлы шаблона, передавая динамические данные из переменной `td`.
	if err != nil {
		app.serverError(wr, err) // Обрабатываем ошибку сервера
	}
}
