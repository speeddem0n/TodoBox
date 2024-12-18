package main

import (
	"html/template"
	"path/filepath"

	"github.com/speeddem0n/todobox/pkg/models"
)

type templateData struct { // Хранилище для любых динамических данных, которые нужно передать HTML-шаблонам.
	Todo  *models.Todo
	Todos []*models.Todo
}

func newTemplateCache(dir string) (map[string]*template.Template, error) { // Функция для кэширования шаблонов

	cache := map[string]*template.Template{} // Инициализируем новую карту, которая будет хранить кэш.

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl")) //Используем функцию filepath.Glob, чтобы получить срез всех файловых путей срасширением '.page.tmpl'.
	if err != nil {
		return nil, err
	}

	for _, page := range pages { // Перебираем файл шаблона от каждой страницы.

		name := filepath.Base(page) // Извлечение конечное названия файла (например, 'home.page.tmpl') из полного пути к файлу и присваивание его переменной name.

		ts, err := template.ParseFiles(page) // Обрабатываем итерируемый файл шаблона.
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl")) // Используем метод ParseGlob для добавления всех каркасных шаблонов. Например base.layout.tmpl
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl")) // Используем метод ParseGlob для добавления всех вспомогательных шаблонов. В нашем случае это footer.partial.tmpl "подвал" нашего шаблона.
		if err != nil {
			return nil, err
		}

		cache[name] = ts // Добавляем полученный набор шаблонов в кэш, используя название страницы
	}

	return cache, nil // Возвращаем полученную карту
}
