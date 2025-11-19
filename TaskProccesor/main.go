package main

import (
	"context"
	"fmt"
	"log"
	"taskProcessor/config"
	"taskProcessor/database"
	"taskProcessor/models"
	"taskProcessor/repository"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde .env (opcional, para desarrollo)
	godotenv.Load()

	// Por ahora, asegúrate de tener las variables en tu sistema o .env cargado

	// 1. Cargar configuración
	cfg := config.Load()

	// 2. Conectar a MongoDB
	mongoDB, err := database.Connect(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		log.Fatalf("❌ Error de conexión: %v", err)
	}
	defer mongoDB.Disconnect()

	// 3. Crear repositorio
	taskRepo := repository.NewTaskRepository(mongoDB.GetCollection("tasks"))

	// 4. Probar operaciones
	ctx := context.Background()

	fmt.Println("\n📝 Probando el repositorio...")

	// === CREAR TAREAS ===
	fmt.Println("\n➕ Creando nuevas tareas...")
	tasks := []*models.Task{
		models.NewTask("Enviar email de bienvenida", map[string]interface{}{
			"email":   "user@example.com",
			"subject": "Bienvenido",
		}),
		models.NewTask("Procesar imagen", map[string]interface{}{
			"image_url": "https://example.com/image.jpg",
			"format":    "thumbnail",
		}),
		models.NewTask("Generar reporte", map[string]interface{}{
			"report_type": "monthly",
			"user_id":     12345,
		}),
	}

	for _, task := range tasks {
		if err := taskRepo.Create(ctx, task); err != nil {
			log.Printf("❌ Error al crear tarea: %v", err)
		} else {
			fmt.Printf("✅ Tarea creada: %s (ID: %s)\n", task.Title, task.ID.Hex())
		}
	}

	// === LISTAR TAREAS PENDIENTES ===
	fmt.Println("\n📋 Listando tareas pendientes...")
	pending, err := taskRepo.FindPending(ctx, 10)
	if err != nil {
		log.Printf("❌ Error: %v", err)
	} else {
		fmt.Printf("   Encontradas: %d tareas\n", len(pending))
		for i, task := range pending {
			fmt.Printf("   %d. %s (Intentos: %d)\n", i+1, task.Title, task.Attempts)
		}
	}

	// === RECLAMAR UNA TAREA ===
	fmt.Println("\n🏷️  Reclamando una tarea como 'worker-test'...")
	claimed, err := taskRepo.ClaimTask(ctx, "worker-test")
	if err != nil {
		log.Printf("❌ Error: %v", err)
	} else if claimed != nil {
		fmt.Printf("✅ Tarea reclamada: %s\n", claimed.Title)
		fmt.Printf("   Worker: %s\n", claimed.ClaimedBy)
		fmt.Printf("   Intentos: %d\n", claimed.Attempts)

		// === MARCAR COMO PROCESADA ===
		fmt.Println("\n✔️  Marcando tarea como procesada...")
		result := fmt.Sprintf("Procesada exitosamente por %s", claimed.ClaimedBy)
		if err := taskRepo.MarkAsProcessed(ctx, claimed.ID, result); err != nil {
			log.Printf("❌ Error: %v", err)
		} else {
			fmt.Println("✅ Tarea procesada exitosamente")
		}
	} else {
		fmt.Println("⚠️  No hay tareas disponibles para reclamar")
	}

	// === ESTADÍSTICAS ===
	fmt.Println("\n📊 Estadísticas finales:")
	total, _ := taskRepo.CountAll(ctx)
	pendingCount, _ := taskRepo.CountPending(ctx)
	fmt.Printf("   Total: %d | Pendientes: %d | Procesadas: %d\n",
		total, pendingCount, total-pendingCount)

	// === LISTAR TODAS LAS TAREAS ===
	fmt.Println("\n📚 Todas las tareas en la base de datos:")
	allTasks, err := taskRepo.FindAll(ctx, 0)
	if err != nil {
		log.Printf("❌ Error: %v", err)
	} else {
		for i, task := range allTasks {
			status := "❌ Pendiente"
			if task.Processed {
				status = "✅ Procesada"
			} else if task.ClaimedBy != "" {
				status = "🔄 En proceso"
			}
			fmt.Printf("   %d. [%s] %s\n", i+1, status, task.Title)
		}
	}

	fmt.Println("\n✨ Prueba completada!")
}
