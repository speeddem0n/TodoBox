package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(wr http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(wr http.ResponseWriter, status int) {
	http.Error(wr, http.StatusText(status), status)
}

func (app *application) notFound(wr http.ResponseWriter) {
	app.clientError(wr, http.StatusNotFound)
}
