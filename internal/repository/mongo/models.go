package mongo

import (
	"time"

	"github.com/turbekoff/todo/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type Session struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Owner    primitive.ObjectID `bson:"owner"`
	Device   string             `bson:"device"`
	Token    string             `bson:"token"`
	ExpireAt time.Time          `bson:"expireAt"`
}

type Task struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Owner     primitive.ObjectID `bson:"owner"`
	Name      string             `bson:"name"`
	Completed bool               `bson:"completed"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

func toUserModel(entity *entities.User) *User {
	id, _ := primitive.ObjectIDFromHex(entity.ID)
	return &User{
		ID:        id,
		Name:      entity.Name,
		Password:  entity.Password,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func toUserEntity(model *User) *entities.User {
	return &entities.User{
		ID:        model.ID.Hex(),
		Name:      model.Name,
		Password:  model.Password,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func toSessionModel(entity *entities.Session) *Session {
	id, _ := primitive.ObjectIDFromHex(entity.ID)
	owner, _ := primitive.ObjectIDFromHex(entity.Owner)
	return &Session{
		ID:       id,
		Owner:    owner,
		Device:   entity.Device,
		Token:    entity.Token,
		ExpireAt: entity.ExpireAt,
	}
}

func toSessionEntity(model *Session) *entities.Session {
	return &entities.Session{
		ID:       model.ID.Hex(),
		Owner:    model.Owner.Hex(),
		Device:   model.Device,
		Token:    model.Token,
		ExpireAt: model.ExpireAt,
	}
}

func toTaskModel(entity *entities.Task) *Task {
	id, _ := primitive.ObjectIDFromHex(entity.ID)
	owner, _ := primitive.ObjectIDFromHex(entity.Owner)
	return &Task{
		ID:        id,
		Owner:     owner,
		Name:      entity.Name,
		Completed: entity.Completed,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func toTaskEntity(entity *Task) *entities.Task {
	return &entities.Task{
		ID:        entity.ID.Hex(),
		Owner:     entity.Owner.Hex(),
		Name:      entity.Name,
		Completed: entity.Completed,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
