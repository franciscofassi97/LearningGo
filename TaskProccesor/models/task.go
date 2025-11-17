package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Title       string              `bson:"title" json:"title"`
	Payload     string              `bson:"payload" json:"payload"`
	Processed   bool                `bson:"processed" json:"processed"`
	Attempts    int                 `bson:"attempts" json:"attempts"`
	ClaimedBy   string              `bson:"claimed_by,omitempty" json:"claimed_by,omitempty"`
	ClaimedAt   *primitive.DateTime `bson:"claimed_at,omitempty" json:"claimed_at,omitempty"`
	ProcessedAt *primitive.DateTime `bson:"processed_at,omitempty" json:"processed_at,omitempty"`
	Result      string              `bson:"result,omitempty" json:"result,omitempty"`
	CreatedAt   time.Time           `bson:"created_at" json:"created_at"`
}

// Creat tarea
func NewTask(title, payload string) *Task {
	return &Task{
		ID:        primitive.NewObjectID(),
		Title:     title,
		Payload:   payload,
		Processed: false,
		Attempts:  0,
		CreatedAt: time.Now(),
	}
}
