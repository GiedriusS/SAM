package api

import (
	"encoding/json"
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
	a := &API{r: r, s: s}
	r.HandleFunc("/hash/{hash}", a.GetAlertByHash)
	r.HandleFunc("/alert/{name}", a.GetRelated)
	return a
}

// GetAlertByHash returns the alert data according to the specified hash.
func (a *API) GetAlertByHash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if data, ok := a.s.Alerts[vars["hash"]]; ok {
		b := []byte{}
		err := json.Unmarshal(b, data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(b))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// GetRelated gets the list of related alerts according to the name and label set.
func (a *API) GetRelated(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "a")
}
