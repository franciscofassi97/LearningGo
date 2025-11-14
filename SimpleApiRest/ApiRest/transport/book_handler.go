package transport

import (
	"apirest/models"
	"apirest/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type BookHandler struct {
	service *service.Service
}

func New(service *service.Service) *BookHandler {
	return &BookHandler{
		service: service,
	}
}

func (handler *BookHandler) HandleBooks(writer http.ResponseWriter, request *http.Request) {
	// Aquí iría la lógica para manejar las solicitudes HTTP relacionadas con los libros
	switch request.Method {
	case http.MethodGet:
		books, err := handler.service.GetAllBooks()
		if err != nil {
			http.Error(writer, "Error al obtener los libros", http.StatusInternalServerError)
			return
		}
		// Aquí se convertirían los libros a JSON y se escribirían en la respuesta
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		json.NewEncoder(writer).Encode(books)

	case http.MethodPost:
		var book models.Book
		if err := json.NewDecoder(request.Body).Decode(&book); err != nil {
			http.Error(writer, "Error al decodificar el libro", http.StatusBadRequest)
			return
		}
		createdBook, err := handler.service.CreateBook(&book)

		if err != nil {
			http.Error(writer, "Error al crear el libro", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		json.NewEncoder(writer).Encode(createdBook)

	default:
		http.Error(writer, "Método no permitido", http.StatusMethodNotAllowed)

	}
}

func (handler *BookHandler) HandleBookByID(writer http.ResponseWriter, request *http.Request) {
	// Aquí iría la lógica para manejar las solicitudes HTTP relacionadas con un libro específico por ID
	idString := strings.TrimPrefix(request.URL.Path, "/books/")
	if idString == "" {
		http.Error(writer, "ID del libro no proporcionado", http.StatusBadRequest)
		return
	}
	// Convertir idString a entero
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(writer, "ID del libro inválido", http.StatusBadRequest)
		return
	}

	switch request.Method {
	case http.MethodGet:
		// Lógica para obtener un libro por ID
		books, err := handler.service.GetBookById(id)
		if err != nil {
			http.Error(writer, "Error al obtener el libro", http.StatusInternalServerError)
			return
		}
		// Aquí se convertiría el libro a JSON y se escribiría en la respuesta'
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(books)

	case http.MethodPut:
		// Lógica para actualizar un libro por ID
		var book models.Book
		if err := json.NewDecoder(request.Body).Decode(&book); err != nil {
			http.Error(writer, "Error al decodificar el libro", http.StatusBadRequest)
			return
		}

		updateBook, err := handler.service.UpdateBook(id, &book)
		if err != nil {
			http.Error(writer, "Error al actualizar el libro", http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(updateBook)

	case http.MethodDelete:
		// Lógica para eliminar un libro por ID
		err := handler.service.DeleteBook(id)
		if err != nil {
			http.Error(writer, "Error al eliminar el libro", http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusNoContent)

	default:
		http.Error(writer, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
