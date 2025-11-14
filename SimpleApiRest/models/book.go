// Package models contiene las estructuras de datos (modelos) de la aplicación
// Es equivalente a los modelos/schemas en Node.js o las interfaces en TypeScript
package models

// Book representa la estructura de un libro en la API
// En Go, los structs son como las clases u objetos en otros lenguajes
type Book struct {
	// ID es el identificador único del libro
	// `json:"id"` es un "tag" que indica cómo se llamará este campo al convertir a JSON
	ID int `json:"id"`

	// Title es el título del libro
	// Los campos que empiezan con mayúscula son públicos (exportados)
	Title string `json:"title"`

	// Author es el autor del libro
	// Si empezaran con minúscula serían privados (solo accesibles dentro del package)
	Author string `json:"author"`
}
