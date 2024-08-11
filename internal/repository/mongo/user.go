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

type UserRepository struct {
	db *mongo.Collection
}

func NewUserRepositry(db *mongo.Database) repositories.UserRepository {
	return &UserRepository{db: db.Collection("users")}
}

func (r *UserRepository) Create(user *entities.User) error {
	model := toUserModel(user)
	_, err := r.db.InsertOne(context.Background(), model)
	return err
}

func (r *UserRepository) Read(id string) (*entities.User, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)

	var user User
	if err := r.db.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repositories.ErrUserNotFound
		}
		return nil, err
	}
	return toUserEntity(&user), nil
}

func (r *UserRepository) ReadByName(name string) (*entities.User, error) {
	var user User
	if err := r.db.FindOne(context.Background(), bson.M{"name": name}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repositories.ErrUserNotFound
		}
		return nil, err
	}
	return toUserEntity(&user), nil
}

func (r *UserRepository) Update(user *entities.User) error {
	model := toUserModel(user)
	query := bson.M{}
	query["name"] = model.Name
	query["password"] = model.Password
	query["createdAt"] = model.CreatedAt
	query["updatedAt"] = model.UpdatedAt

	_, err := r.db.UpdateOne(context.Background(), bson.M{"_id": model.ID}, bson.M{"$set": query})
	return err
}

func (r *UserRepository) Delete(id string) error {
	objectID, _ := primitive.ObjectIDFromHex(id)

	_, err := r.db.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}
