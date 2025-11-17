# Task Processor

Una aplicación en Go para encolar y procesar tareas concurrentemente usando MongoDB.

## Características

- API REST para crear tareas
- Procesamiento concurrente con worker pool
- MongoDB para persistencia
- Procesamiento atómico sin duplicados

## Requisitos

- Go 1.24+
- MongoDB (local o Atlas)

## Configuración

1. Copia `.env.example` a `.env`:

   ```bash
   cp .env.example .env
   ```

2. Edita `.env` con tu conexión de MongoDB:

   ```env
   MONGODB_URI=mongodb+srv://usuario:password@cluster.mongodb.net/
   MONGODB_DATABASE=taskProcessor
   ```

3. Instala las dependencias:
   ```bash
   go mod download
   ```

## Uso

```bash
go run main.go
```

## Estructura del Proyecto

```
TaskProccesor/
├── main.go              # Punto de entrada
├── go.mod               # Dependencias
├── .env                 # Configuración (no subir a git)
└── .env.example         # Ejemplo de configuración
```

## Pasos de Desarrollo

- [x] Paso 0: Setup inicial del proyecto
- [ ] Paso 1: Estructura de carpetas + conexión MongoDB
- [ ] Paso 2: Modelo Task + repository básico
- [ ] Paso 3: Servicio con ProcessTasks y worker pool
- [ ] Paso 4: Handlers HTTP y server
- [ ] Paso 5: Pruebas unitarias
- [ ] Paso 6: Mejoras y optimizaciones
