package mongo

import (
	"context"
	"errors"

	"github.com/turbekoff/todo/internal/domain/entities"
	"github.com/turbekoff/todo/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository struct {
	db *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) repositories.TaskRepository {
	return &TaskRepository{db: db.Collection("tasks")}
}

func (r *TaskRepository) Create(task *entities.Task) error {
	model := toTaskModel(task)
	_, err := r.db.InsertOne(context.Background(), model)
	return err
}

func (r *TaskRepository) Read(id string) (*entities.Task, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)

	var task Task
	if err := r.db.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&task); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repositories.ErrTaskNotFound
		}
		return nil, err
	}
	return toTaskEntity(&task), nil
}

func (r *TaskRepository) ReadAllByOwner(owner string) ([]*entities.Task, error) {
	objectID, _ := primitive.ObjectIDFromHex(owner)

	cursor, err := r.db.Find(context.Background(), bson.M{"owner": objectID})
	if err != nil {
		return nil, err
	}

	var tasks []Task
	if err = cursor.All(context.TODO(), &tasks); err != nil {
		return nil, err
	}

	var entities []*entities.Task
	for _, task := range tasks {
		entities = append(entities, toTaskEntity(&task))
	}

	return entities, nil
}

func (r *TaskRepository) Update(task *entities.Task) error {
	model := toTaskModel(task)
	query := bson.M{}
	query["owner"] = model.Owner
	query["name"] = model.Name
	query["completed"] = model.Completed
	query["createdAt"] = model.CreatedAt
	query["updatedAt"] = model.UpdatedAt

	_, err := r.db.UpdateOne(context.Background(), bson.M{"_id": model.ID}, bson.M{"$set": query})
	return err
}

func (r *TaskRepository) Delete(id string) error {
	objectID, _ := primitive.ObjectIDFromHex(id)

	_, err := r.db.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}
