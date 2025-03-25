package main

import (
	"GOST_update/config"
	"GOST_update/db"
	"GOST_update/parser"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func downloadFile(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Пример использования: go run main.go \"ГОСТ Р 1.1-2020\"")
		return
	}

	gostNumber := strings.Join(os.Args[1:], " ")

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// Подключаемся к базе данных в начале
	pg, err := db.NewPostgres(cfg.DBURL)
	if err != nil {
		fmt.Println("Error connecting to DB:", err)
		return
	}

	filePath := "katalog-natsionalinyh-standartov-06-08-2024.csv"
	needParse := false

	// Проверяем, существует ли CSV-файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Если файл не существует, загружаем его
		url := "https://www.rst.gov.ru/opendata/7706406291-nationalstandards/data-20240808-structure-20220330.csv"
		fmt.Println("CSV-файл не найден. Загружаем...")
		err = downloadFile(url, filePath)
		if err != nil {
			fmt.Println("Ошибка загрузки файла:", err)
			return
		}
		fmt.Println("CSV-файл успешно загружен.")
		needParse = true
	} else {
		// Если файл существует, проверяем, есть ли данные в базе
		isEmpty, err := pg.IsTableEmpty()
		if err != nil {
			fmt.Println("Ошибка при проверке базы данных:", err)
			return
		}
		needParse = isEmpty
	}

	// Парсим CSV только если нужно
	if needParse {
		fmt.Println("Начинаем парсинг CSV-файла...")
		records, err := parser.ParseCSV(filePath)
		if err != nil {
			fmt.Println("Error parsing CSV:", err)
			return
		}

		err = pg.SaveRecords(records)
		if err != nil {
			fmt.Println("Error saving records:", err)
			return
		}
		fmt.Println("Данные успешно загружены в базу данных!")
	}

	// Проверяем, есть ли введённый стандарт в базе данных
	exists, err := pg.CheckRecordExists(gostNumber)
	if err != nil {
		fmt.Println("Ошибка проверки записи:", err)
		return
	}

	// Выводим результат
	if exists {
		fmt.Printf("'%s' есть в системе.\n", gostNumber)
	} else {
		fmt.Printf("'%s' нет в системе.\n", gostNumber)
	}
}
