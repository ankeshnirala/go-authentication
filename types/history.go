package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogUserHistory struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId    primitive.ObjectID `json:"userId"`
	LogType   string             `json:"logType"`
	CreatedAt string             `json:"createdAt"`
}

func NewUserLog(userId primitive.ObjectID, logType string) (*LogUserHistory, error) {
	return &LogUserHistory{
		UserId:    userId,
		LogType:   logType,
		CreatedAt: time.Now().UTC().String(),
	}, nil
}
