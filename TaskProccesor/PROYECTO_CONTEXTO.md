# Task Processor - Contexto del Proyecto

## 📋 DESCRIPCIÓN GENERAL

Estamos desarrollando una aplicación en Go llamada **"Task Processor"** que funciona como una cola de tareas con procesamiento concurrente usando MongoDB como base de datos.

### Características principales:

- API REST para crear y gestionar tareas
- Procesamiento concurrente seguro con worker pool
- MongoDB para persistencia de datos
- Operaciones atómicas para evitar procesamiento duplicado de tareas
- Uso de `primitive.ObjectID` para IDs de MongoDB
- Timeouts con context en todas las operaciones de base de datos

---

## 🎯 ENFOQUE DE DESARROLLO

Estamos siguiendo un **enfoque iterativo paso a paso**:

### Por cada paso:

1. Solo implementamos el código necesario para ese paso específico
2. Creamos archivos específicos y organizados
3. Explicamos claramente qué hace cada parte
4. Probamos antes de continuar al siguiente paso

### Pasos planificados:

- ✅ **Paso 0**: Setup inicial del proyecto
- ✅ **Paso 1**: Estructura de carpetas + conexión MongoDB
- ✅ **Paso 2**: Modelo Task + repository básico
- ⏳ **Paso 3**: Servicio con ProcessTasks y worker pool (SIGUIENTE)
- 📅 **Paso 4**: Handlers HTTP y server
- 📅 **Paso 5**: Pruebas unitarias
- 📅 **Paso 6**: Mejoras y optimizaciones

---

## 📊 MODELO DE DATOS

### Estructura de Task en MongoDB:

```go
type Task struct {
    ID          primitive.ObjectID     // _id único de MongoDB
    Title       string                 // Título descriptivo de la tarea
    Payload     map[string]interface{} // Datos arbitrarios de la tarea
    Processed   bool                   // ¿Ya fue procesada?
    Attempts    int                    // Número de intentos de procesamiento
    ClaimedBy   string                 // ID del worker que la reclamó
    ClaimedAt   *time.Time            // Timestamp cuando fue reclamada
    ProcessedAt *time.Time            // Timestamp cuando fue procesada
    Result      string                 // Resultado del procesamiento
    CreatedAt   time.Time             // Timestamp de creación
}
```

### Campos importantes:

- **`ClaimedBy` y `ClaimedAt`**: Permiten rastrear qué worker está procesando la tarea
- **`Attempts`**: Se incrementa cada vez que un worker reclama la tarea
- **`Processed`**: Flag booleano para filtrar tareas completadas

---

## 🏗️ ARQUITECTURA ACTUAL

### Estructura de carpetas:

```
TaskProccesor/
├── main.go                      # Punto de entrada
├── go.mod                       # Dependencias
├── go.sum                       # Checksums
├── .env                         # Variables de entorno (no subir a git)
├── .env.example                 # Ejemplo de configuración
├── .gitignore                   # Archivos ignorados
├── README.md                    # Documentación
├── config/
│   └── config.go                # Carga configuración desde env vars
├── database/
│   └── mongodb.go               # Conexión a MongoDB
├── models/
│   └── task.go                  # Definición del modelo Task
└── repository/
    └── task_repository.go       # Operaciones CRUD con MongoDB
```

### Capas de la aplicación:

1. **Config**: Gestiona variables de entorno
2. **Database**: Maneja conexión a MongoDB
3. **Models**: Define estructuras de datos
4. **Repository**: Acceso a datos (CRUD)
5. **Service**: Lógica de negocio (próximo paso)
6. **Handlers**: Endpoints HTTP (paso posterior)

---

## 🔧 CONFIGURACIÓN

### Variables de entorno (.env):

```env
MONGODB_URI=mongodb+srv://usuario:password@cluster.mongodb.net/
MONGODB_DATABASE=taskProcessor
SERVER_PORT=8080
WORKER_COUNT=5
```

### Dependencias (go.mod):

```
go.mongodb.org/mongo-driver v1.17.6
github.com/joho/godotenv
```

---

## 💾 REPOSITORY - OPERACIONES IMPLEMENTADAS

### Métodos del TaskRepository:

#### 1. **Create(ctx, task)**

Crea una nueva tarea en MongoDB.

#### 2. **FindByID(ctx, id)**

Busca una tarea por su ObjectID.

#### 3. **FindAll(ctx, limit)**

Lista todas las tareas ordenadas por fecha (descendente).

#### 4. **FindPending(ctx, limit)**

Busca tareas no procesadas y no reclamadas.

```go
filter := bson.M{
    "processed":  false,
    "claimed_by": bson.M{"$exists": false},
}
```

#### 5. **ClaimTask(ctx, workerID)** ⭐ **OPERACIÓN CLAVE**

Reclama una tarea atómicamente usando `FindOneAndUpdate`.

- **Atómica**: Previene que múltiples workers reclamen la misma tarea
- Incrementa `attempts` automáticamente
- Asigna `claimed_by` y `claimed_at`
- Retorna la tarea actualizada o `nil` si no hay disponibles

```go
filter := bson.M{
    "processed":  false,
    "claimed_by": bson.M{"$exists": false},
}
update := bson.M{
    "$set": bson.M{
        "claimed_by": workerID,
        "claimed_at": now,
    },
    "$inc": bson.M{"attempts": 1},
}
```

#### 6. **MarkAsProcessed(ctx, id, result)**

Marca una tarea como completada con su resultado.

#### 7. **CountAll(ctx)** y **CountPending(ctx)**

Retornan estadísticas de tareas.

---

## 🔒 SEGURIDAD EN CONCURRENCIA

### ¿Cómo evitamos duplicados?

1. **Operación atómica**: `FindOneAndUpdate` ejecuta "buscar + actualizar" en una sola operación
2. **Filtro específico**: Solo busca tareas sin `claimed_by`
3. **Primera coincidencia**: MongoDB garantiza que solo un worker obtiene la tarea

### Flujo de procesamiento seguro:

```
Worker 1 y Worker 2 ejecutan ClaimTask() simultáneamente
    ↓
MongoDB procesa FindOneAndUpdate atómicamente
    ↓
Worker 1 obtiene Task A (claimed_by = "worker-1")
Worker 2 obtiene Task B (claimed_by = "worker-2")
    ↓
No hay duplicados ✅
```

---

## 🧪 PRUEBAS REALIZADAS

### Test en main.go (Paso 2):

1. ✅ Conectar a MongoDB exitosamente
2. ✅ Crear 3 tareas de ejemplo
3. ✅ Listar tareas pendientes
4. ✅ Reclamar una tarea con `ClaimTask`
5. ✅ Marcar tarea como procesada
6. ✅ Mostrar estadísticas (total, pendientes, procesadas)

### Resultado esperado:

```
🚀 Iniciando Task Processor...
✅ Conectado exitosamente a MongoDB
📝 Probando el repositorio...

➕ Creando nuevas tareas...
✅ Tarea creada: Enviar email de bienvenida
✅ Tarea creada: Procesar imagen
✅ Tarea creada: Generar reporte

📋 Tareas pendientes: 3

🏷️  Reclamando una tarea...
✅ Tarea reclamada: Enviar email de bienvenida por worker-test

✔️  Marcando tarea como procesada...
✅ Tarea procesada exitosamente

📊 Estadísticas:
  Total: 3 | Pendientes: 2 | Procesadas: 1
```

---

## 📝 DECISIONES DE IMPLEMENTACIÓN

### 1. **¿Por qué primitive.ObjectID?**

- Es el tipo nativo de MongoDB para IDs únicos
- Incluye timestamp de creación
- Más eficiente que strings para indexación

### 2. **¿Por qué Context con timeout?**

```go
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
```

- Evita operaciones bloqueadas indefinidamente
- Permite cancelación de operaciones
- Mejora la resiliencia de la aplicación

### 3. **¿Por qué separar Repository y Service?**

- **Repository**: Solo acceso a datos (CRUD puro)
- **Service**: Lógica de negocio + orquestación de workers
- Facilita testing y mantenimiento

### 4. **¿Por qué map[string]interface{} para Payload?**

- Permite datos arbitrarios sin esquema fijo
- Flexible para diferentes tipos de tareas
- Se serializa naturalmente a BSON

---

## 🎯 PRÓXIMO PASO: Worker Pool

### Paso 3: Servicio con ProcessTasks y worker pool

Lo que implementaremos:

1. **`service/task_service.go`**:

   - `StartWorkerPool(workerCount)`: Iniciar N workers concurrentes
   - `ProcessTasks()`: Buscar y procesar tareas continuamente
   - `processTask(task)`: Lógica de procesamiento de una tarea individual
   - `StopWorkers()`: Detener workers gracefully

2. **Conceptos a usar**:

   - **Goroutines**: Para ejecutar workers concurrentemente
   - **Channels**: Para comunicación entre workers
   - **WaitGroup**: Para esperar que todos los workers terminen
   - **Context**: Para cancelación coordinada

3. **Flujo esperado**:
   ```
   main.go
     ↓
   StartWorkerPool(5) → Lanza 5 goroutines
     ↓
   Cada worker:
     - ClaimTask() atómicamente
     - Procesa la tarea (simular trabajo)
     - MarkAsProcessed()
     - Repite
   ```

---

## 📚 CÓDIGO DE REFERENCIA

### Ejemplo de creación de tarea:

```go
task := models.NewTask("Enviar email", map[string]interface{}{
    "email": "user@example.com",
    "subject": "Bienvenido",
})

err := taskRepo.Create(context.Background(), task)
```

### Ejemplo de procesamiento atómico:

```go
// Worker reclama tarea
claimed, err := taskRepo.ClaimTask(ctx, "worker-1")
if claimed != nil {
    // Procesar...
    result := processTask(claimed)

    // Marcar como procesada
    taskRepo.MarkAsProcessed(ctx, claimed.ID, result)
}
```

---

## 🔗 INFORMACIÓN ADICIONAL

### Base de datos MongoDB:

- Nombre: `taskProcessor`
- Colección: `tasks`
- Conexión: MongoDB Atlas (cloud)

### Entorno de desarrollo:

- Go version: 1.24.2
- OS: Windows
- Shell: PowerShell

---

## ✅ ESTADO ACTUAL

### Completado:

- [x] Setup inicial del proyecto
- [x] Sistema de configuración con variables de entorno
- [x] Conexión a MongoDB con manejo de errores
- [x] Modelo Task con todos los campos necesarios
- [x] Repository completo con operaciones CRUD
- [x] Operación atómica `ClaimTask` para concurrencia segura
- [x] Pruebas básicas de funcionalidad

### Pendiente:

- [ ] Servicio con worker pool
- [ ] API REST con handlers HTTP
- [ ] Pruebas unitarias
- [ ] Optimizaciones y mejoras

---

## 💡 CONCEPTOS CLAVE PARA ENTENDER

1. **Procesamiento atómico**: Garantiza que una tarea solo sea procesada por un worker
2. **Context timeout**: Todas las operaciones de DB tienen límite de tiempo
3. **Worker pool**: Grupo de goroutines que procesan tareas concurrentemente
4. **BSON**: Formato binario de MongoDB (similar a JSON)
5. **Defer**: Garantiza ejecución de cleanup (ej: `defer cancel()`)

---

## 🎓 ENFOQUE EDUCATIVO

Este proyecto está diseñado para aprender:

- Estructuración de proyectos Go
- Trabajo con MongoDB en Go
- Concurrencia segura con goroutines
- Manejo de contextos y timeouts
- Patrones de repository y service
- APIs REST con Go

---

**Última actualización**: Paso 2 completado
**Siguiente paso**: Implementar worker pool con goroutines
