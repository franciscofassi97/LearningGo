package repository

import (
	"context"
	"fmt"
	"taskProcessor/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(collection *mongo.Collection) *TaskRepository {
	return &TaskRepository{
		collection: collection,
	}
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if task.ID.IsZero() {
		task.ID = primitive.NewObjectID()
	}

	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}

	_, err := r.collection.InsertOne(ctx, task)
	if err != nil {
		return fmt.Errorf("error al crear trarea: %v", err)
	}
	return nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var task models.Task
	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error al obtener tarea por ID: %v", err)
	}

	return &task, nil
}

func (r *TaskRepository) FindAll(ctx context.Context, limit int64) ([]*models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	if limit > 0 {
		opts.SetLimit(limit)
	}

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("error al listar tareas:  %v", err)
	}
	defer cursor.Close(ctx)

	var tasks []*models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, fmt.Errorf("error al decodificar tareas: %v", err)
	}

	return tasks, nil
}

func (r *TaskRepository) FindPending(ctx context.Context, limit int64) ([]*models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{
		"processed":  false,
		"claimed_by": bson.M{"$exists": false},
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})
	if limit > 0 {
		opts.SetLimit(limit)
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error al listar tareas pendientes: %v", err)
	}
	defer cursor.Close(ctx)

	var task []*models.Task
	if err = cursor.All(ctx, &task); err != nil {
		return nil, fmt.Errorf("error al decodificar tareas pendientes: %v", err)
	}
	return task, nil
}

func (r *TaskRepository) ClaimTask(ctx context.Context, workerID string) (*models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	now := time.Now()

	filter := bson.M{
		"processed":  false,
		"claimed_by": bson.M{"$exists": false},
	}

	update := bson.M{
		"$set": bson.M{
			"claimed_by": workerID,
			"claimed_at": now, // Ahora es time.Time
		},
		"$inc": bson.M{"attempts": 1},
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetSort(bson.D{{Key: "created_at", Value: 1}})

	var task models.Task
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error al reclamar tarea: %v", err)
	}
	return &task, nil
}

func (r *TaskRepository) MarkAsProcessed(ctx context.Context, id primitive.ObjectID, result string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	now := time.Now()
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"processed":    true,
			"processed_at": now, // Ahora es time.Time
			"result":       result,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error al marcar tarea como procesada: %v", err)
	}
	return nil
}

func (r *TaskRepository) CountAll(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("error al contar tareas: %v", err)
	}
	return count, nil
}

func (r *TaskRepository) CountPending(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	filter := bson.M{
		"processed":  false,
		"claimed_by": bson.M{"$exists": false},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("error al contar tareas pendientes: %v", err)
	}
	return count, nil
}
