package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Passowrd string `json:"password"`
}

type LoginResponse struct {
	UserID primitive.ObjectID `json:"userID"`
	Token  string             `json:"token"`
}

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	InsertedID primitive.ObjectID `json:"insertedId"`
	Token      string             `json:"token"`
}

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
	CreatedOn string             `json:"createdOn"`
	UpdatedOn string             `json:"updatedOn"`
}

func (u *User) ValidPassword(pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd)) == nil
}

func NewUser(name, email, password string) (*User, error) {

	encPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Name:      name,
		Email:     email,
		Password:  string(encPwd),
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	}, nil
}
