package handlers

import (
	"net/http"
)

func handleErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "Something Went Wrong")
}

