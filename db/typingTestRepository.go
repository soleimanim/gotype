package db

import (
	"time"

	"github.com/jmoiron/sqlx"
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
