package service

import (
	"apirest/models"
	"apirest/store"
)

type Service struct {
	store store.Store
}

func New(store store.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetAllBooks() ([]*models.Book, error) {
	return s.store.GetAll()
}

func (s *Service) GetBookById(id int) (*models.Book, error) {
	return s.store.GetById(id)
}

func (s *Service) CreateBook(book *models.Book) (*models.Book, error) {
	return s.store.Create(book)
}

func (s *Service) UpdateBook(id int, book *models.Book) (*models.Book, error) {
	return s.store.Update(id, book)
}

func (s *Service) DeleteBook(id int) error {
	return s.store.Delete(id)
}
