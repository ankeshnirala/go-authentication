package storage

import (
	"github.com/ankeshnirala/go/authentication/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStorage interface {
	RegisterUser(*types.User) (*mongo.InsertOneResult, error)
	GetUserByID(id primitive.ObjectID) *mongo.SingleResult
	GetUserByEmail(string) *mongo.SingleResult
	LogUserHistory(*types.LogUserHistory) (*mongo.InsertOneResult, error)
	GetLogsByUserID(id primitive.ObjectID) (*mongo.Cursor, error)
}
