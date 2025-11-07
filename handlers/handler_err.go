package handlers

import (
	"net/http"

	"github.com/khaingminhtun/rssagg/utilis"
)

func HandlerError(w http.ResponseWriter, r *http.Request) {
	utilis.RespondWithError(w, 400, "somethind went wrong")
}
