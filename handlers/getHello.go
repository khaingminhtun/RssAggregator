package handlers

import (
	"net/http"

	"github.com/khaingminhtun/rssagg/utilis"
)

func GetHello(w http.ResponseWriter, r *http.Request) {
	data := struct{}{}
	utilis.RespondWithJSON(w, http.StatusOK, data)
}
