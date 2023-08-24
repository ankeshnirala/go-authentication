package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ankeshnirala/go/authentication/utils"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StandardResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// here we are making apiFunc as our handler so that we pass it's type
// to makeHTTPHandleFunc for handling the error
type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

// we have written this function mux handle func does not return any error
// so we need to handle that returned error manually
func MakeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle the error
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	res := &StandardResponse{
		Code: status,
		Data: v,
	}

	return json.NewEncoder(w).Encode(res)
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Reading token from headers
		tokenString := r.Header.Get("x-jwt-token")

		cookie, err := r.Cookie("name")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(cookie)

		// checking there is a token or not
		token, err := utils.ValidateJWT(tokenString)
		if err != nil {
			PermissionDenied(w, err)
			return
		}

		if !token.Valid {
			PermissionDenied(w, nil)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			PermissionDenied(w, nil)
			return
		}

		userID, err := primitive.ObjectIDFromHex(claims["userID"].(string))
		if err != nil {
			PermissionDenied(w, err)
			return
		}

		// 3) check if user still exist
		// var user *types.User
		// s.GetUserByID(userId).Decode(&user)

		// if user == nil {
		// 	permissionDenied(w)
		// 	return
		// }

		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
