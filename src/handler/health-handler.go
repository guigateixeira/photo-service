package handler

import (
	"net/http"

	"photo-service/src/util"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	util.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "running", "success": "true"})
}
