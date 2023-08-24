package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ankeshnirala/go/authentication/constants"
	"github.com/ankeshnirala/go/authentication/types"
	"github.com/ankeshnirala/go/authentication/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Server) SignupHandler(w http.ResponseWriter, r *http.Request) error {
	// check if method is POST or not
	if r.Method != "POST" {
		return fmt.Errorf(constants.METHODNOTALLOWED, r.Method)
	}

	// sync request body data with SignupRequest
	req := new(types.SignupRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	// check if user already exist or not
	var existingUser *types.User
	s.mongoStore.GetUserByEmail(req.Email).Decode(&existingUser)
	if existingUser != nil {
		return fmt.Errorf(constants.EMAILREGISTERED, req.Email)
	}

	// create a new user
	user, err := types.NewUser(req.Name, req.Email, req.Password)
	if err != nil {
		return err
	}

	// register new user in db
	result, err := s.mongoStore.RegisterUser(user)
	if err != nil {
		return err
	}
	s.logger.Printf(constants.USERREGISTERED, result.InsertedID)

	// generate jwt token and send it in response to login
	token, err := utils.CreateJWT(result.InsertedID.(primitive.ObjectID))
	if err != nil {
		return err
	}
	s.logger.Printf(constants.TOKENCREATED, result.InsertedID)

	// store user logs history
	ulogs, err := types.NewUserLog(result.InsertedID.(primitive.ObjectID), "registration")
	if err != nil {
		return err
	}

	log, err := s.mongoStore.LogUserHistory(ulogs)
	if err != nil {
		return err
	}
	s.logger.Printf(constants.LOGADDED, log.InsertedID)

	// set token in http cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: utils.ExpirationTime,
	})

	// configure signup response
	res := types.SignupResponse{
		InsertedID: result.InsertedID.(primitive.ObjectID),
		Token:      token,
	}

	return WriteJSON(w, http.StatusOK, res)
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	// check if method is POST or not
	if r.Method != "POST" {
		return fmt.Errorf(constants.METHODNOTALLOWED, r.Method)
	}

	// sync request body data with LoginRequest
	req := new(types.LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	// check if user already exist or not
	var userExist *types.User
	s.mongoStore.GetUserByEmail(req.Email).Decode(&userExist)
	if userExist == nil {
		return fmt.Errorf(constants.EMAILNOTREGISTERED, req.Email)
	}

	// check password
	if ok := userExist.ValidPassword(req.Passowrd); !ok {
		return fmt.Errorf(constants.INCORRECTPWS)
	}

	// generate jwt token and send it in response to login
	token, err := utils.CreateJWT(userExist.ID)
	if err != nil {
		return err
	}
	s.logger.Printf(constants.TOKENCREATED, userExist.ID)

	// store user logs history
	ulogs, err := types.NewUserLog(userExist.ID, "login")
	if err != nil {
		return err
	}

	log, err := s.mongoStore.LogUserHistory(ulogs)
	if err != nil {
		return err
	}
	s.logger.Printf(constants.LOGADDED, log.InsertedID)

	// set token in http cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: utils.ExpirationTime,
	})

	// configure signup response
	res := types.LoginResponse{
		UserID: userExist.ID,
		Token:  token,
	}

	return WriteJSON(w, http.StatusOK, res)
}

func (s *Server) LogHistoryHandler(w http.ResponseWriter, r *http.Request) error {
	// check if method is POST or not
	if r.Method != "GET" {
		return fmt.Errorf(constants.METHODNOTALLOWED, r.Method)
	}

	userId := r.Context().Value("userID")

	userLogs, err := s.mongoStore.GetLogsByUserID(userId.(primitive.ObjectID))
	if err != nil {
		return err
	}

	var usrLogs []*types.LogUserHistory
	if err := userLogs.All(context.Background(), &usrLogs); err != nil {
		return err
	}

	// Reading token from headers
	// tokenString := r.Header.Get("x-jwt-token")

	return WriteJSON(w, http.StatusOK, usrLogs)
}
