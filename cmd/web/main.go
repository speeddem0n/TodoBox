package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	pSQL "github.com/speeddem0n/todobox/pkg/models/postgreSQL"

	_ "github.com/lib/pq"
)

type application struct { //Зависимости для наших обработчиков
	errorLog      *log.Logger                   // Логер для ошибок
	infoLog       *log.Logger                   // Логер для INFO
	todos         *pSQL.TodoModel               // postgreSQL model
	templateCache map[string]*template.Template // map с кэшом html шаблонов
}

func main() {

	addr := flag.String("addr", ":4000", "Network addres HTTP")                                                                // flag для адрела сервера
	dsn := flag.String("dsn", "user=web password=1235 dbname=todobox sslmode=disable", "Название postgreSQL источника данных") // флаг для подключения к БД

	flag.Parse() // Получаем флаг из командной строки

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                  // Логер для записи информационных сообщений
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // Логер для записи ошибок

	db, err := openDB(*dsn) // Подключене в Базе данных
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close() // Откладываем закрытие подключения к БД по завершению работы приложения

	templateCache, err := newTemplateCache(".\\ui\\html\\") // Инициализируем новый кэш шаблона
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{ // создаем новую структура со всеми зависимостями
		errorLog:      errorLog,
		infoLog:       infoLog,
		todos:         &pSQL.TodoModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{ // Структура для подключения к серверу
		Addr:     *addr,        // Server addres
		ErrorLog: errorLog,     // Custom logger
		Handler:  app.routes(), // Mux Handler
	}

	infoLog.Printf("Запуск веб-сервера на %s", *addr) // Сообщение о том что сервер запущен по адресу ...
	err = srv.ListenAndServe()                        // Запускаем сервер
	errorLog.Fatal(err)                               // Обарабатываем ошибку
}

func openDB(dsn string) (*sql.DB, error) { // Функия для открытия подключения к БД
	db, err := sql.Open("postgres", dsn) // Инициализируем подключение к БД
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil { // Ping подтверждает наличие подключения к БД
		return nil, err
	}

	return db, nil
}
