package main

import (
	"apirest/service"
	"apirest/store"
	"apirest/transport"
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func main() {
	//Conectar a sql lite
	database, err := sql.Open("sqlite", "./books.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	//Crear la tabla si no existe.
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL
	);`
	_, err = database.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	//Inyectar nuestras dependencias
	bookStore := store.New(database)
	bookService := service.New(bookStore)
	bookHandle := transport.New(bookService)

	//Configurar el router y los endpoints
	http.HandleFunc("/books", bookHandle.HandleBooks)
	http.HandleFunc("/books/", bookHandle.HandleBookByID)

	//Iniciar el servidor
	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
