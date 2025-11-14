# Apuntes del Proyecto - API REST en Go

## 📚 Lo que Aprendimos Hoy

### 1. Inicialización de Proyectos Go
```bash
go mod init nombre-del-proyecto
```
- Crea `go.mod` (equivalente a `package.json` en Node.js)
- Define el módulo y gestiona dependencias
- `go.sum` se genera automáticamente (como `package-lock.json`)

### 2. Compilación vs Interpretación
**Go es compilado:**
- `go run main.go` → Compila en memoria y ejecuta
- `go build` → Genera un `.exe` (binario ejecutable)
- El `.exe` es autónomo (no necesitas Go instalado para ejecutarlo)

**Diferencia con Node.js:**
- Node: Interpretado, necesitas Node instalado
- Go: Compilado a binario nativo, distribuyes solo el `.exe`

### 3. Punteros (`*`)
```go
func Create(book *models.Book) (*models.Book, error)
```

**¿Qué es?** Referencia a la dirección de memoria del dato

**¿Por qué usarlo?**
- **Eficiencia**: Pasa 8 bytes (dirección) en lugar de copiar toda la estructura
- **Modificación**: Permite modificar el valor original
- **Nulabilidad**: Puede ser `nil`

**Importante:** Al convertir a JSON, Go automáticamente "desreferencia" el puntero:
```go
book := &models.Book{ID: 1, Title: "Don Quijote"}
json.Marshal(book)  // {"id":1,"title":"Don Quijote"}
```

### 4. Import con `_` (Blank Identifier)
```go
import _ "modernc.org/sqlite"
```

**¿Para qué?** Importa el paquete SOLO por sus efectos secundarios (side effects)

**Uso común:** Drivers de base de datos que se auto-registran:
```go
// El driver internamente hace:
func init() {
    sql.Register("sqlite", &SQLiteDriver{})
}
```

No necesitas llamar funciones del paquete directamente.

### 5. Slices (`[]`)
```go
var books []*models.Book  // Slice de punteros a Book
books = append(books, &book)  // Agrega elementos dinámicamente
```

**Diferencia con arrays:**
- Arrays: Tamaño fijo → `var books [10]Book`
- Slices: Dinámicos → `var books []Book`

### 6. Manejo de Errores - Múltiples Retornos
```go
func GetAll() ([]*models.Book, error)

// Uso:
books, err := store.GetAll()
if err != nil {
    return nil, err
}
```

**Go NO tiene try/catch** - Los errores se manejan explícitamente.

### 7. Structs y Tags JSON
```go
type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}
```

**Tags:** Metadatos que indican cómo serializar/deserializar JSON

**Visibilidad:**
- **Mayúscula** (`Title`): Público/Exportado
- **minúscula** (`title`): Privado del paquete

### 8. Interfaces
```go
type Store interface {
    GetAll() ([]*models.Book, error)
}
```

**Implementación implícita:** Si un tipo tiene todos los métodos de la interface, automáticamente la implementa.

No necesitas `implements` como en Java/TypeScript.

### 9. Receivers (Métodos)
```go
func (s *store) GetAll() ([]*models.Book, error) { ... }
```

**`(s *store)`** = receiver (como `this` en otros lenguajes)

**Tipos:**
- `(s *store)` - Receiver por puntero (puede modificar s)
- `(s store)` - Receiver por valor (copia)

### 10. `defer`
```go
rows, err := db.Query(...)
defer rows.Close()  // Se ejecuta al final de la función
```

Equivalente a `finally` - Se ejecuta SIEMPRE antes de retornar.

---

## 🗺️ Rutas HTTP

### Definición en `main.go`
```go
http.HandleFunc("/books", bookHandle.HandleBooks)      // Exactamente /books
http.HandleFunc("/books/", bookHandle.HandleBookByID)  // /books/* (con ID)
```

### Extracción del ID
```go
idString := strings.TrimPrefix(request.URL.Path, "/books/")
// /books/1 → "1"
id, err := strconv.Atoi(idString)  // String a int
```

### Routing por Método HTTP
```go
switch request.Method {
case http.MethodGet:
    // Lógica GET
case http.MethodPost:
    // Lógica POST
default:
    http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
}
```

---

## 🗄️ SQLite en Go

### Instalación
```bash
# Opción 1: Requiere GCC
go get github.com/mattn/go-sqlite3

# Opción 2: Pure Go (sin GCC) - LA QUE USAMOS
go get modernc.org/sqlite
```

### Conexión
```go
database, err := sql.Open("sqlite", "./books.db")
defer database.Close()
```

El archivo `.db` se crea automáticamente si no existe.

### Queries

**Múltiples filas:**
```go
rows, err := db.Query("SELECT id, title, author FROM books")
defer rows.Close()

for rows.Next() {
    rows.Scan(&book.ID, &book.Title, &book.Author)
    books = append(books, &book)
}
```

**Una fila:**
```go
err := db.QueryRow("SELECT ... WHERE id = ?", id).Scan(&book.ID, ...)
```

**Modificar:**
```go
result, err := db.Exec("INSERT INTO books (title, author) VALUES (?, ?)", title, author)
id, _ := result.LastInsertId()
```

**`?` = placeholders** que previenen SQL injection.

---

## 🏗️ Arquitectura del Proyecto

```
SimpleApiRest/
├── models/       → Estructuras de datos (DTOs)
├── store/        → Acceso a datos (Repository)
├── service/      → Lógica de negocio
├── transport/    → Handlers HTTP (Controllers)
└── main.go       → Punto de entrada
```

**Patrón:** Clean Architecture / Layered Architecture

**Beneficio:** Cada capa es independiente y testeable.

---

## 🔄 Flujo de una Request

**Ejemplo: `POST /books`**

```
Cliente
  ↓ HTTP POST /books + JSON
main.go (Servidor)
  ↓ Busca handler
transport/book_handler.go
  ↓ HandleBooks() - Detecta POST
  ↓ Decodifica JSON a struct
service/book_services.go
  ↓ CreateBook() - Lógica de negocio
store/book_store.go
  ↓ Create() - INSERT en SQLite
  ↓ Retorna libro con ID
Cliente
  ← {"id":1,"title":"...","author":"..."}
```

---

## 🔌 Inyección de Dependencias

```go
bookStore := store.New(database)        // Store depende de DB
bookService := service.New(bookStore)   // Service depende de Store
bookHandle := transport.New(bookService) // Handler depende de Service
```

**Ventajas:**
- Testeable (inyectas mocks)
- Flexible (cambias implementaciones)
- Dependencias explícitas

---

## 📡 Endpoints de la API

| Método | Ruta | Descripción | Body |
|--------|------|-------------|------|
| `GET` | `/books` | Listar todos | - |
| `GET` | `/books/1` | Obtener por ID | - |
| `POST` | `/books` | Crear libro | `{"title":"...","author":"..."}` |
| `PUT` | `/books/1` | Actualizar | `{"title":"...","author":"..."}` |
| `DELETE` | `/books/1` | Eliminar | - |

---

## 💻 Comandos PowerShell para Probar

### Crear libro
```powershell
Invoke-WebRequest -Uri http://localhost:8080/books -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"title":"Cien años de soledad","author":"Gabriel García Márquez"}'
```

### Listar todos
```powershell
Invoke-WebRequest -Uri http://localhost:8080/books
```

### Obtener por ID
```powershell
Invoke-WebRequest -Uri http://localhost:8080/books/1
```

### Actualizar
```powershell
Invoke-WebRequest -Uri http://localhost:8080/books/1 -Method PUT -Headers @{"Content-Type"="application/json"} -Body '{"title":"Don Quijote","author":"Cervantes"}'
```

### Eliminar
```powershell
Invoke-WebRequest -Uri http://localhost:8080/books/1 -Method DELETE
```

---

## 🐛 Problemas Comunes y Soluciones

### Error: "nil dereference"
**Problema:**
```go
var book *models.Book  // nil
book.Title = "..."     // ❌ Error!
```

**Solución:**
```go
book := &models.Book{}  // Inicializa
book.Title = "..."      // ✅ OK
```

### Error: Ruta no funciona con ID
**Problema:**
```go
idString := strings.TrimPrefix(request.URL.Path, "/book/")  // ❌ Singular
```

**Solución:**
```go
idString := strings.TrimPrefix(request.URL.Path, "/books/")  // ✅ Plural
```

Debe coincidir con la ruta registrada en `main.go`.

---

## 🎯 Conceptos Clave de Go

### No hay constructores
Se usan funciones `New()`:
```go
func New(db *sql.DB) Store {
    return &store{db: db}
}
```

### Retornar interfaces, no tipos concretos
```go
func New(db *sql.DB) Store {  // Interface, no *store
    return &store{db: db}
}
```

### `&` obtiene la dirección de memoria
```go
&store{db: db}  // Retorna puntero al struct
```

### Patrón Repository
`store/` abstrae acceso a datos - puedes cambiar de SQLite a PostgreSQL sin modificar otras capas.

---

## 📦 Git & GitHub

### Configuración inicial
```bash
git config --global user.name "Tu Nombre"
git config --global user.email "tu@email.com"
```

### Crear repositorio
```bash
git init
git add .
git commit -m "Initial commit"
```

### Subir a GitHub
```bash
git remote add origin https://github.com/usuario/repo.git
git branch -M main
git push -u origin main
```

### `.gitignore` para Go
```
*.exe
*.db
*.test
*.out
go.work
.vscode/
vendor/
```

---

## 🚀 Ventajas de Go

1. **Binario compilado** - Un `.exe` autónomo
2. **Rápido** - Muy eficiente
3. **Concurrencia nativa** - Goroutines
4. **Tipado estático** - Menos errores
5. **Batería incluida** - HTTP, JSON, SQL integrados
6. **Sin runtime** - No necesitas Go instalado para ejecutar

---

## 📚 Recursos de Aprendizaje

- [Tour of Go](https://go.dev/tour/) - Tutorial oficial
- [Effective Go](https://go.dev/doc/effective_go) - Mejores prácticas
- [Go by Example](https://gobyexample.com/) - Ejemplos prácticos
- [Documentación estándar](https://pkg.go.dev/std)

---

## ✅ Checklist de lo que Hicimos

- [x] Inicializar módulo Go (`go mod init`)
- [x] Crear estructura de proyecto (models, store, service, transport)
- [x] Implementar CRUD completo
- [x] Configurar SQLite (driver pure Go)
- [x] Crear handlers HTTP
- [x] Manejar rutas y métodos HTTP
- [x] Probar API con PowerShell
- [x] Crear README técnico
- [x] Configurar Git y `.gitignore`
- [x] Subir a GitHub
- [x] Aprender conceptos fundamentales de Go

---

## 🎓 Próximos Pasos

1. **Agregar validaciones** en el service layer
2. **Implementar middleware** para logging
3. **Usar Gin framework** para routing más avanzado
4. **Agregar tests unitarios**
5. **Implementar autenticación JWT**
6. **Dockerizar la aplicación**
7. **Agregar documentación con Swagger**

---

**Fecha:** 14 de Noviembre, 2025  
**Proyecto:** Simple REST API en Go con SQLite  
**Repositorio:** https://github.com/franciscofassi97/LearningGo
