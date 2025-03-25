package parser

import (
	"GOST_update/models"
	"encoding/csv"
	"os"

	"golang.org/x/text/encoding/charmap" // Пакет для работы с кодировками
	"golang.org/x/text/transform"        // Пакет для преобразования данных
)

func ParseCSV(filePath string) ([]models.Record, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Преобразователь из Windows-1251 в UTF-8
	decoder := charmap.Windows1251.NewDecoder()
	reader := transform.NewReader(file, decoder)

	// Читаем CSV-файл
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []models.Record
	for i, row := range records {
		if i == 0 {
			continue
		}

		data = append(data, models.Record{
			Number: row[0],
			Name:   row[1],
			State:  row[2],
		})
	}

	return data, nil
}
