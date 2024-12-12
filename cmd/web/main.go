package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct { //Зависимости для наших обработчиков
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	addr := flag.String("addr", ":4000", "Network addres HTTP") // Флаг командной строки, значение по умолчанию: ":4000"

	flag.Parse() // Получаем флаг из командной строки

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                  //Логер для записи информационных сообщений
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) //Логер для записи ошибок

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *addr,        // Server addres
		ErrorLog: errorLog,     // Custom logger
		Handler:  app.routes(), // Mux Handler
	}

	infoLog.Printf("Запуск веб-сервера на %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
