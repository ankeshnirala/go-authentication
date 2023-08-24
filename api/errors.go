package api

import "net/http"

func PermissionDenied(w http.ResponseWriter, err error) {
	if err != nil {
		WriteJSON(w, http.StatusForbidden, ApiError{Error: err.Error()})
		return
	}

	WriteJSON(w, http.StatusForbidden, ApiError{Error: "Permission denied"})
}
