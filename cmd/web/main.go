package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"

	pSQL "github.com/speeddem0n/todobox/pkg/models/postgreSQL"

	_ "github.com/lib/pq"
)

type application struct { //Зависимости для наших обработчиков
	errorLog      *log.Logger                   // Логер для ошибок
	infoLog       *log.Logger                   // Логер для INFO
	todos         *pSQL.TodoModel               // postgreSQL model
	templateCache map[string]*template.Template // map с кэшом html шаблонов
}

type config struct { // Структура для хранения настроек сервера
	Addr       string `yaml:"addr"`       // Адрес сервера
	UsernameDB string `yaml:"usernameDB"` // Имя пользователя БД
	PasswordDB string `yaml:"passwordDB"` // Пароль пользователя БД
	DBname     string `yaml:"DBname"`     // Название базы данных
	Sslmode    string `yaml:"sslmode"`    // SSLmode enable or disable
}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                  // Логер для записи информационных сообщений
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // Логер для записи ошибок

	conf := config{} // Создаем пустую структуру для записи настроек из файла

	file, err := os.ReadFile("config.yaml") // Считываем файл с настройками
	if err != nil {
		errorLog.Fatal(err) // Приложение прекращает работу в случае ошибки
	}

	err = yaml.Unmarshal(file, &conf) // Декадируем YAML файл в структуру config
	if err != nil {
		errorLog.Fatal(err) // Приложение прекращает работу в случае ошибки
	}

	db, err := openDB(&conf) // Подключене в Базе данных
	if err != nil {
		errorLog.Fatal(err) // Приложение прекращает работу в случае ошибки
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
		Addr:     *&conf.Addr,  // Server addres
		ErrorLog: errorLog,     // Custom logger
		Handler:  app.routes(), // Mux Handler
	}

	infoLog.Printf("Запуск веб-сервера на %s", *&conf.Addr) // Сообщение о том что сервер запущен по адресу ...
	err = srv.ListenAndServe()                              // Запускаем сервер
	errorLog.Fatal(err)                                     // Обарабатываем ошибку
}

func openDB(conf *config) (*sql.DB, error) { // Функия для открытия подключения к БД
	psqlInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", // Преобразуем параметры из структуры config в одну строку для  sql.Open
		conf.UsernameDB, conf.PasswordDB, conf.DBname, conf.Sslmode)
	db, err := sql.Open("postgres", psqlInfo) // Инициализируем подключение к БД
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil { // Ping подтверждает наличие подключения к БД
		return nil, err
	}

	return db, nil
}
