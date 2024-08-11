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

type SessionRepository struct {
	db *mongo.Collection
}

func NewSessionRepository(db *mongo.Database) repositories.SessionRepository {
	return &SessionRepository{db: db.Collection("sessions")}
}

func (r *SessionRepository) Create(session *entities.Session) error {
	model := toSessionModel(session)
	_, err := r.db.InsertOne(context.Background(), model)
	return err
}

func (r *SessionRepository) Read(id string) (*entities.Session, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)

	var session Session
	if err := r.db.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&session); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repositories.ErrSessionNotFound
		}
		return nil, err
	}
	return toSessionEntity(&session), nil
}

func (r *SessionRepository) ReadAllByOwner(owner string) ([]*entities.Session, error) {
	objectID, _ := primitive.ObjectIDFromHex(owner)

	cursor, err := r.db.Find(context.Background(), bson.M{"owner": objectID})
	if err != nil {
		return nil, err
	}

	var sessions []Session
	if err = cursor.All(context.TODO(), &sessions); err != nil {
		return nil, err
	}

	var entities []*entities.Session
	for _, session := range sessions {
		entities = append(entities, toSessionEntity(&session))
	}

	return entities, nil
}

func (r *SessionRepository) Update(session *entities.Session) error {
	model := toSessionModel(session)
	query := bson.M{}
	query["owner"] = model.Owner
	query["device"] = model.Device
	query["token"] = model.Token
	query["expireAt"] = model.ExpireAt

	_, err := r.db.UpdateOne(context.Background(), bson.M{"_id": model.ID}, bson.M{"$set": query})
	return err
}

func (r *SessionRepository) Delete(id string) error {
	objectID, _ := primitive.ObjectIDFromHex(id)

	_, err := r.db.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}
