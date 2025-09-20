package models

import (
	"database/sql"
	"time"
)

type Classifier struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}

type ClassifierModel struct {
	DB *sql.DB
}

func (m *ClassifierModel) Insert(name string) (int64, error) {
	query := `INSERT INTO classifiers (name) VALUES (?)`
	result, err := m.DB.Exec(query, name)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *ClassifierModel) Get(id int64) (*Classifier, error) {
	query := `SELECT id, name, created_at FROM classifiers WHERE id = ?`
	row := m.DB.QueryRow(query, id)
	var c Classifier
	err := row.Scan(&c.ID, &c.Name, &c.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (m *ClassifierModel) List() ([]*Classifier, error) {
	query := `SELECT id, name, created_at FROM classifiers ORDER BY created_at DESC`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classifiers []*Classifier
	for rows.Next() {
		var c Classifier
		err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		classifiers = append(classifiers, &c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return classifiers, nil
}
