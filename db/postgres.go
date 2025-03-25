package db

import (
	"GOST_update/models"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(dbURL string) (*Postgres, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}
	return &Postgres{db: db}, nil
}

func (p *Postgres) SaveRecords(records []models.Record) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO records (number, name, state) VALUES ($1, $2, $3)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, record := range records {
		_, err = stmt.Exec(record.Number, record.Name, record.State)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (p *Postgres) CheckRecordExists(gostNumber string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM records WHERE number = $1)`
	var exists bool
	err := p.db.QueryRow(query, gostNumber).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	return exists, nil
}

func (p *Postgres) IsTableEmpty() (bool, error) {
	var count int
	err := p.db.QueryRow("SELECT COUNT(*) FROM records").Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
