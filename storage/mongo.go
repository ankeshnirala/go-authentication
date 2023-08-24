package storage

import (
	"context"
	"os"
	"time"

	"github.com/ankeshnirala/go/authentication/constants"
	"github.com/ankeshnirala/go/authentication/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type MongoStore struct {
	db *mongo.Database
}

func NewMongoStore() (*MongoStore, error) {
	connStr := os.Getenv("MONGODB_URL")
	DATABASE := os.Getenv("DATABASE_NAME")

	clientOptions := options.Client().ApplyURI(connStr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &MongoStore{db: client.Database(DATABASE)}, nil
}

func (s *MongoStore) GetUserByID(id primitive.ObjectID) *mongo.SingleResult {
	return s.db.Collection(constants.USERS).FindOne(context.TODO(), bson.M{"_id": id})
}

func (s *MongoStore) GetUserByEmail(email string) *mongo.SingleResult {
	return s.db.Collection(constants.USERS).FindOne(context.TODO(), bson.M{"email": email})
}

func (s *MongoStore) RegisterUser(u *types.User) (*mongo.InsertOneResult, error) {
	return s.db.Collection(constants.USERS).InsertOne(context.TODO(), u)
}

func (s *MongoStore) LogUserHistory(u *types.LogUserHistory) (*mongo.InsertOneResult, error) {
	return s.db.Collection(constants.USERLOGS).InsertOne(context.TODO(), u)
}

func (s *MongoStore) GetLogsByUserID(id primitive.ObjectID) (*mongo.Cursor, error) {
	return s.db.Collection(constants.USERLOGS).Find(context.TODO(), bson.M{"userid": id})
}
