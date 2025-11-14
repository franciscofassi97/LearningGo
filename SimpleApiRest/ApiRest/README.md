# API REST de Libros en Go - Guía Técnica

Este proyecto es una API REST para gestionar libros, construida en Go con SQLite. Esta guía explica **cómo funciona** y **por qué** se tomaron ciertas decisiones técnicas.

---

## 📂 Arquitectura del Proyecto

```
ApiRest/
├── models/          # Estructuras de datos
├── store/           # Capa de acceso a datos (Repository)
├── service/         # Lógica de negocio
├── transport/       # Handlers HTTP (Controllers)
├── main.go          # Punto de entrada
├── go.mod           # Definición del módulo y dependencias
└── books.db         # Base de datos SQLite (se genera automáticamente)
```

### ¿Por qué esta arquitectura?

Esta estructura sigue el patrón **Clean Architecture** / **Layered Architecture**:

- **models**: Define las estructuras de datos (DTOs/Entities)
- **store**: Abstrae el acceso a la base de datos (Repository Pattern)
- **service**: Contiene la lógica de negocio (Business Logic)
- **transport**: Maneja HTTP y convierte requests/responses (Presentation Layer)

**Beneficio**: Cada capa tiene una responsabilidad única y puede ser modificada/probada independientemente.

---

## 🔑 Conceptos Clave de Go

### 1. **Punteros (`*`)**

```go
func Create(book *models.Book) (*models.Book, error)
```

**¿Qué es `*`?**  
Es un **puntero** - una referencia a la dirección de memoria donde está almacenado el dato.

**¿Por qué usar punteros?**
- **Eficiencia**: Pasas la dirección de memoria (8 bytes) en lugar de copiar toda la estructura
- **Modificación**: Permite modificar el valor original
- **Nulabilidad**: Un puntero puede ser `nil` (null)

**Ejemplo práctico:**
```go
// Sin puntero - se copia toda la estructura
func UpdateBook(book models.Book) {
    book.Title = "Nuevo título"  // ❌ Modifica la copia, no el original
}

// Con puntero - se pasa la referencia
func UpdateBook(book *models.Book) {
    book.Title = "Nuevo título"  // ✅ Modifica el original
}
```

**En JSON no cambia nada:**
```go
book := &models.Book{ID: 1, Title: "Don Quijote"}
json.Marshal(book)  // {"id":1,"title":"Don Quijote"}
```
Go automáticamente "desreferencia" el puntero al convertir a JSON.

---

### 2. **Import con Guion Bajo (`_`)**

```go
import _ "modernc.org/sqlite"
```

**¿Qué significa `_`?**  
Se llama **blank identifier** y significa: "Importa este paquete SOLO por sus efectos secundarios, no voy a usar sus funciones directamente".

**¿Por qué se usa con drivers de BD?**  
Los drivers de base de datos se auto-registran al importarse:

```go
// Dentro de modernc.org/sqlite hay algo así:
func init() {
    sql.Register("sqlite", &SQLiteDriver{})  // Se registra automáticamente
}
```

No necesitas llamar funciones del paquete directamente, solo que se ejecute su `init()`.

---

### 3. **Slices (`[]`)**

```go
func GetAll() ([]*models.Book, error)
```

**¿Qué es `[]`?**  
Es un **slice** - similar a un array dinámico (como arrays en JavaScript o listas en Python).

```go
var books []*models.Book  // Slice de punteros a Book
books = append(books, &book)  // Agrega elementos dinámicamente
```

**Diferencia con arrays:**
- Arrays tienen tamaño fijo: `var books [10]Book`
- Slices son dinámicos: `var books []Book`

---

### 4. **Múltiples Valores de Retorno**

```go
func GetAll() ([]*models.Book, error)
```

Go permite retornar **múltiples valores**. Es el patrón estándar para manejo de errores:

```go
books, err := store.GetAll()
if err != nil {
    // Manejar error
    return nil, err
}
// Usar books
```

**No hay try/catch en Go** - los errores se manejan explícitamente con este patrón.

---

### 5. **Structs y Tags JSON**

```go
type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}
```

**Structs**: Son como clases/objetos en otros lenguajes.

**Tags (`json:"campo"`)**: Metadatos que indican cómo serializar/deserializar JSON:
```go
book := Book{ID: 1, Title: "Don Quijote", Author: "Cervantes"}
json.Marshal(book)
// Resultado: {"id":1,"title":"Don Quijote","author":"Cervantes"}
```

**Visibilidad:**
- **Mayúscula**: Público/Exportado (`Title` - accesible desde otros paquetes)
- **minúscula**: Privado (`title` - solo accesible dentro del paquete)

---

### 6. **Interfaces**

```go
type Store interface {
    GetAll() ([]*models.Book, error)
    GetById(id int) (*models.Book, error)
    // ...
}
```

**¿Qué es una interface?**  
Define un **contrato** - un conjunto de métodos que un tipo debe implementar.

**Implementación implícita:**
```go
type store struct {
    db *sql.DB
}

// Si store tiene todos los métodos de Store, automáticamente implementa Store
func (s *store) GetAll() ([]*models.Book, error) { ... }
```

No necesitas escribir `implements Store` como en Java/TypeScript.

**Beneficio**: Puedes cambiar la implementación (SQLite → PostgreSQL) sin modificar el código que usa la interface.

---

### 7. **Receivers (Métodos)**

```go
func (s *store) GetAll() ([]*models.Book, error) { ... }
```

**¿Qué es `(s *store)`?**  
Es el **receiver** - similar a `this` en otros lenguajes. Indica que este método pertenece al tipo `*store`.

```go
bookStore := store.New(db)
bookStore.GetAll()  // s dentro del método es bookStore
```

**Puntero vs Valor:**
- `(s *store)` - Receiver por puntero (puede modificar s, más eficiente)
- `(s store)` - Receiver por valor (recibe una copia)

---

### 8. **`defer`**

```go
defer rows.Close()
```

**¿Qué hace `defer`?**  
Ejecuta la función **al final** de la función actual, sin importar cómo termina (normal o con error).

**Equivalente a `finally` en otros lenguajes:**
```go
func GetAll() ([]*models.Book, error) {
    rows, err := db.Query(...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()  // Se ejecuta SIEMPRE antes de retornar
    
    // ... procesar rows ...
    return books, nil
}  // rows.Close() se ejecuta aquí
```

---

## 🗺️ Manejo de Rutas HTTP

### Definición de Rutas (`main.go`)

```go
http.HandleFunc("/books", bookHandle.HandleBooks)
http.HandleFunc("/books/", bookHandle.HandleBookByID)
```

**¿Cómo funciona?**

1. **`/books`** (sin `/` final):
   - Coincide **exactamente** con `/books`
   - `GET /books` → `HandleBooks`
   - `POST /books` → `HandleBooks`

2. **`/books/`** (con `/` final):
   - Coincide con `/books/` **y todo lo que siga**
   - `GET /books/1` → `HandleBookByID`
   - `DELETE /books/42` → `HandleBookByID`

### Extracción del ID

```go
idString := strings.TrimPrefix(request.URL.Path, "/books/")
// /books/1 → "1"
// /books/42 → "42"

id, err := strconv.Atoi(idString)  // Convierte string a int
```

### Enrutamiento por Método HTTP

```go
func HandleBooks(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        // Listar todos los libros
    case http.MethodPost:
        // Crear libro
    default:
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
    }
}
```

Go no tiene un router sofisticado por defecto, por eso usamos `switch` para manejar diferentes métodos HTTP.

---

## 🔄 Flujo de una Request

**Ejemplo: `POST /books` (Crear libro)**

```
1. Cliente → HTTP POST /books
              Body: {"title":"Don Quijote","author":"Cervantes"}

2. main.go → Servidor recibe la request
              Busca handler para /books
              
3. transport/book_handler.go → HandleBooks()
              Detecta método POST
              Decodifica JSON a struct Book
              
4. service/book_services.go → CreateBook()
              Valida datos (si hubiera lógica de negocio)
              
5. store/book_store.go → Create()
              Ejecuta INSERT en SQLite
              Obtiene ID auto-generado
              
6. Respuesta ← {"id":1,"title":"Don Quijote","author":"Cervantes"}
```

---

## 🗄️ Base de Datos SQLite

### Conexión

```go
database, err := sql.Open("sqlite", "./books.db")
```

- **`sql.Open`**: Función del paquete estándar `database/sql`
- **`"sqlite"`**: Nombre del driver (registrado por `modernc.org/sqlite`)
- **`"./books.db"`**: Ruta del archivo de base de datos (se crea si no existe)

### Creación de Tabla

```go
createTableQuery := `
CREATE TABLE IF NOT EXISTS books (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    author TEXT NOT NULL
);`
database.Exec(createTableQuery)
```

- **`IF NOT EXISTS`**: Solo crea la tabla si no existe (idempotente)
- **`AUTOINCREMENT`**: SQLite genera automáticamente IDs únicos

### Queries

**Query múltiples filas:**
```go
rows, err := db.Query("SELECT id, title, author FROM books")
for rows.Next() {
    rows.Scan(&book.ID, &book.Title, &book.Author)
}
```

**Query una sola fila:**
```go
err := db.QueryRow("SELECT ... WHERE id = ?", id).Scan(&book.ID, ...)
```

**Modificar datos:**
```go
result, err := db.Exec("INSERT INTO books (title, author) VALUES (?, ?)", title, author)
id, _ := result.LastInsertId()
```

**`?` son placeholders** - evitan SQL injection al escapar automáticamente los valores.

---

## 🔌 Inyección de Dependencias

```go
// En main.go
bookStore := store.New(database)        // Store depende de DB
bookService := service.New(bookStore)   // Service depende de Store
bookHandle := transport.New(bookService) // Handler depende de Service
```

**¿Por qué?**
- **Testeable**: Puedes inyectar mocks en tests
- **Flexible**: Cambiar implementaciones sin modificar código
- **Claro**: Dependencias explícitas, no ocultas

---

## 📦 Gestión de Módulos

### `go.mod`
```go
module apirest

go 1.25.4
```

- Define el nombre del módulo (`apirest`)
- Especifica la versión de Go
- Lista dependencias (se agregan automáticamente con `go get`)

### Imports
```go
import "apirest/models"
```

Los imports se basan en el nombre del módulo definido en `go.mod`.

---

## 🎯 Endpoints de la API

| Método | Ruta | Descripción | Body |
|--------|------|-------------|------|
| `GET` | `/books` | Listar todos los libros | - |
| `GET` | `/books/:id` | Obtener libro por ID | - |
| `POST` | `/books` | Crear nuevo libro | `{"title":"...","author":"..."}` |
| `PUT` | `/books/:id` | Actualizar libro | `{"title":"...","author":"..."}` |
| `DELETE` | `/books/:id` | Eliminar libro | - |

---

## 🧠 Conceptos Importantes

### 1. **¿Por qué no hay constructores?**
Go no tiene constructores. Por convención se usan funciones `New()`:

```go
func New(db *sql.DB) Store {
    return &store{db: db}
}
```

### 2. **¿Por qué retornar interfaces?**
```go
func New(db *sql.DB) Store {  // Retorna interface, no *store
```
Esto permite cambiar la implementación sin modificar el código que la usa.

### 3. **¿Por qué `&` en `&store{}`?**
```go
return &store{db: db}  // Retorna un puntero al struct
```
`&` obtiene la dirección de memoria (crea un puntero).

### 4. **¿Qué es el patrón Repository?**
`store/` actúa como Repository - abstrae el acceso a datos. Si cambias de SQLite a PostgreSQL, solo modificas `store/`, el resto del código no cambia.

---

## 🚀 Ventajas de Go para APIs

1. **Compilado a binario**: Un solo `.exe`, fácil de desplegar
2. **Concurrencia nativa**: Goroutines para manejar miles de requests
3. **Tipado estático**: Menos errores en runtime
4. **Rápido**: Muy eficiente en memoria y CPU
5. **Batería incluida**: HTTP, JSON, SQL en la biblioteca estándar
6. **Sin runtime externo**: No necesitas Node.js/Python instalado

---

## 📚 Recursos para Aprender Más

- [Tour of Go](https://go.dev/tour/) - Tutorial oficial interactivo
- [Effective Go](https://go.dev/doc/effective_go) - Mejores prácticas
- [Go by Example](https://gobyexample.com/) - Ejemplos prácticos
- [Standard Library](https://pkg.go.dev/std) - Documentación de paquetes estándar
