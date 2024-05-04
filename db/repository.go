package db

type Repository[Model any] interface {
	Create(*Model) error
	GetAll(limit int, offset int) ([]Model, error)
	CountAllWhere(query string) (int, error)
	MaxWhere(field string, query string) (any, error)
	Sum(field string, query string) (any, error)
	Avg(field string, where string) (any, error)
}
