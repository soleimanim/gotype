package db

type Repository[Model any] interface {
	Create(*Model) error
	GetAll(limit int, offset int) ([]Model, error)
}
