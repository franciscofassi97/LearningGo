package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	DataBase *mongo.Database
}

// Connect Establecer conexión a MongoDB
func Connect(uri, databaseName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Configurar opciones de cliente
	clientOptions := options.Client().ApplyURI(uri)

	//conectar al cliente
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, fmt.Errorf("error al conectar a MongoDB: %v", err)
	}

	//verificar conexion con ping
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("error al hacer ping a MongoDB: %v", err)
	}

	log.Printf("✅ Conectado exitosamente a MongoDB (base de datos: %s)", databaseName)

	return &MongoDB{
		Client:   client,
		DataBase: client.Database(databaseName),
	}, nil
}

// Disconnect Desconectar de MongoDB
func (mongo *MongoDB) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongo.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("error al desconectar de MongoDB: %v", err)
	}

	log.Println("✅ Desconectado de MongoDB")

	return nil
}

// GetCollection retorna una coleccion de base de datos
func (mongo *MongoDB) GetCollection(collectionName string) *mongo.Collection {
	return mongo.DataBase.Collection(collectionName)
}
