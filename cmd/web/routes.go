package main

import (
	"net/http"
	"path/filepath"
)

type neuteredFileSystem struct { // Настраиваемая имплементация файловой системы http.FileSystem, с помощью которой будет возвращаться ошибка os.ErrNotExist для любого HTTP запроса напрямую к папке.
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) { // Метод Open(), который вызывается каждый раз, когда http.FileServer получает запрос.
	f, err := nfs.fs.Open(path) // Мы открываем вызываемый путь
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() { // Используя метод IsDir() мы проверим если вызываемый путь является папкой или нет.
		index := filepath.Join(path, "index.html")    // Если это папка, то с помощью метода Stat("index.html") мы проверим если файл index.html существует внутри данной папки
		if _, err := nfs.fs.Open(index); err != nil { // Если файл index.html не существует, то метод вернет ошибку os.ErrNotExist
			closeErr := f.Close() // Метод Close() для закрытия только, что открытого index.html файла, чтобы избежать утечки файлового дескриптора.
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

func (app *application) routes() *http.ServeMux { // Мультиплексор HTTP-запросов
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)             // Домашная страница
	mux.HandleFunc("/todo", app.showTodo)     // Отображения заметки по id Пример "/todo?id=1"
	mux.HandleFunc("/create", app.newTodo)    // Созадние заметки
	mux.HandleFunc("/delete", app.deleteTodo) // Удаление заметки
	mux.HandleFunc("/search", app.searchTodo) // Поиск заметки

	fileServer := http.FileServer(neuteredFileSystem{http.Dir(".\\ui\\static")}) // FileServer возвращает обработчик который обслуживает HTTP-запросы с содержимым файловой системы.

	mux.Handle("/static", http.NotFoundHandler())                   // Защита, что бы не было видно файлом static приложения по конечной ссылке /static
	mux.Handle("/static/", http.StripPrefix("/static", fileServer)) // Требуется удалить "/static" из URL перед его отправкой в http.FileServer с помощью http.StripPrefix(). Иначе пользователь получит ошибку 404

	return mux
}
