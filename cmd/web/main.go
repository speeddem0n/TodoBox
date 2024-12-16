package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	pSQL "github.com/speeddem0n/todobox/pkg/models/postgreSQL"

	_ "github.com/lib/pq"
)

type application struct { //Зависимости для наших обработчиков
	errorLog *log.Logger
	infoLog  *log.Logger
	todos    *pSQL.TodoModel // postgreSQL model
}

func main() {

	addr := flag.String("addr", ":4000", "Network addres HTTP")
	dsn := flag.String("dsn", "user=web password=1235 dbname=todobox sslmode=disable", "Название postgreSQL источника данных")

	flag.Parse() // Получаем флаг из командной строки

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                  //Логер для записи информационных сообщений
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) //Логер для записи ошибок

	db, err := openDB(*dsn) // Подключене в Базе данных
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		todos:    &pSQL.TodoModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,        // Server addres
		ErrorLog: errorLog,     // Custom logger
		Handler:  app.routes(), // Mux Handler
	}

	infoLog.Printf("Запуск веб-сервера на %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
