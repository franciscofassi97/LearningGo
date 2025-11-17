package main

import (
	"log"
	"taskProcessor/config"
	"taskProcessor/database"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde .env (opcional, para desarrollo)
	godotenv.Load()

	// Por ahora, asegúrate de tener las variables en tu sistema o .env cargado

	// Cargar configuración
	cfg := config.Load()
	log.Printf("🚀 Iniciando Task Processor...")
	log.Printf("📊 Base de datos: %s", cfg.MongoDatabase)

	// Conectar a MongoDB
	db, err := database.Connect(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		log.Fatalf("❌ Error fatal: %v", err)
	}
	defer db.Disconnect()

	log.Println("✅ Task Processor está listo")
	log.Println("📝 Presiona Ctrl+C para salir")

	// Mantener el programa corriendo (en próximos pasos será el servidor HTTP)
	select {}
}
