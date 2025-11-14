package store

import (
	"apirest/models"
	"database/sql"
)

type Store interface {
	GetAll() ([]*models.Book, error)
	GetById(id int) (*models.Book, error)
	Create(book *models.Book) (*models.Book, error)
	Update(id int, book *models.Book) (*models.Book, error)
	Delete(id int) error
}

type store struct {
	db *sql.DB
}

func New(db *sql.DB) Store {
	return &store{
		db: db,
	}
}

func (store *store) GetAll() ([]*models.Book, error) {
	query := "SELECT id, title, author FROM books"

	rows, err := store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.Book

	for rows.Next() {
		var book models.Book

		if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
			return nil, err
		}

		books = append(books, &book)
	}
	return books, nil
}

func (store *store) GetById(id int) (*models.Book, error) {
	query := "SELECT id, title, author FROM books WHERE id = ?"
	book := models.Book{}

	err := store.db.QueryRow(query, id).Scan(&book.ID, &book.Title, &book.Author)

	if err != nil {
		return nil, err
	}
	return &book, nil

}

func (store *store) Create(book *models.Book) (*models.Book, error) {
	query := "INSERT INTO books (title, author) VALUES (?, ?)"

	resp, err := store.db.Exec(query, book.Title, book.Author)

	if err != nil {
		return nil, err
	}

	id, err := resp.LastInsertId()

	if err != nil {
		return nil, err
	}

	book.ID = int(id)

	return book, nil
}

func (store *store) Update(id int, book *models.Book) (*models.Book, error) {
	query := "UPDATE books SET title = ?, author = ? WHERE id = ?"
	_, err := store.db.Exec(query, book.Title, book.Author, id)

	if err != nil {
		return nil, err
	}
	return book, nil
}

func (store *store) Delete(id int) error {
	query := "DELETE FROM books WHERE id = ?"

	_, err := store.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
