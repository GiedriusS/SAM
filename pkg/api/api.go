package api

import (
	"fmt"
	"net/http"

	"github.com/GiedriusS/SAM/pkg/alerts"
	"github.com/gorilla/mux"
)

// API is a wrapper around a HTTP router and supporting data.
type API struct {
	r *mux.Router
	s *alerts.State
}

// NewAPI creates a new API object.
func NewAPI(s *alerts.State) *API {
	r := mux.NewRouter()
	r.HandleFunc("/alert/{hash}", GetAlertByHash)
	return &API{r: r, s: s}
}

// GetAlertByHash returns the alert data according to the specified hash.
func GetAlertByHash(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}
