package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/soleimanim/gotype/logger"
)

type TypingTestModel struct {
	ID            int       `db:"id"`
	TestDate      time.Time `db:"test_date"`
	Speed         float32   `db:"speed"`
	Accuracy      float32   `db:"accuracy"`
	WordsCount    uint      `db:"words_count"`
	MistakesCount uint      `db:"mistakes_count"`
}
type TypingTestRepository struct {
	db *sqlx.DB
}

func NewTypingTestRepository(db *sqlx.DB) Repository[TypingTestModel] {
	return &TypingTestRepository{
		db: db,
	}
}

func (r *TypingTestRepository) Create(m *TypingTestModel) error {
	_, err := r.db.Exec("INSERT INTO "+TABLE_TYPING_TESTS+" (speed, accuracy, words_count, mistakes_count) VALUES (?, ?, ?, ?)",
		m.Speed, m.Accuracy, m.WordsCount, m.MistakesCount)
	return err
}

func (r *TypingTestRepository) GetAll(limit int, offset int) ([]TypingTestModel, error) {
	rows, err := r.db.Queryx("SELECT * FROM "+TABLE_TYPING_TESTS+" ORDER BY id DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return []TypingTestModel{}, err
	}
	result := make([]TypingTestModel, 0)
	for rows.Next() {
		var m TypingTestModel
		err := rows.StructScan(&m)
		if err != nil {
			return []TypingTestModel{}, err
		}
		result = append(result, m)
	}

	return result, nil
}

func (r *TypingTestRepository) CountAllWhere(conditions string) (int, error) {
	where := ""
	if conditions != "" {
		where = "where " + conditions
	}
	rows, err := r.db.Query("SELECT COUNT(*) FROM " + TABLE_TYPING_TESTS + where)
	if err != nil {
		logger.Println("Error reading count from database", err)
		return 0, err
	}
	if rows.Next() {
		var count int
		rows.Scan(&count)
		return count, nil
	}

	return 0, errors.New("could not read count from database")
}
func (r *TypingTestRepository) MaxWhere(field string, query string) (any, error) {
	where := ""
	if query != "" {
		where = "WHERE " + query
	}
	rows, err := r.db.Query(fmt.Sprintf("SELECT MAX(%s) FROM %s %s", field, TABLE_TYPING_TESTS, where))
	if err != nil {
		logger.Println("Error finding max value", err)
		return 0, err
	}

	if rows.Next() {
		var maxValue any
		rows.Scan(&maxValue)
		return maxValue, nil
	}
	return 0, errors.New("could not fetch database query result")
}

func (r *TypingTestRepository) Sum(field string, query string) (any, error) {
	where := ""
	if query != "" {
		where = "WHERE " + query
	}
	rows, err := r.db.Query(fmt.Sprintf("SELECT SUM(%s) FROM %s %s", field, TABLE_TYPING_TESTS, where))
	if err != nil {
		logger.Println("Error reading sum from database", err)
		return 0, err
	}
	if rows.Next() {
		var sum any
		rows.Scan(&sum)
		return sum, nil
	}
	return 0, errors.New("could not fetch query result")
}

func (r TypingTestRepository) Avg(field string, conditions string) (any, error) {
	where := ""
	if conditions != "" {
		where = "WHERE " + conditions
	}
	q := fmt.Sprintf("SELECT AVG(%s) FROM %s %s", field, TABLE_TYPING_TESTS, where)
	rows, err := r.db.Query(q)
	if err != nil {
		logger.Println("Error on typing test repository Avg", err)
		return 0, err
	}

	if rows.Next() {
		var val any
		rows.Scan(&val)
		return val, nil
	}

	return 0, errors.New("could not fetch query result")
}
