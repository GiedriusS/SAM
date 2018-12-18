package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GiedriusS/SAM/pkg/alerts"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// API is a wrapper around a HTTP router and supporting data.
type API struct {
	R *mux.Router
	s *alerts.State
	l *zap.Logger
}

// NewAPI creates a new API object.
func NewAPI(s *alerts.State, logger *zap.Logger) *API {
	r := mux.NewRouter()
	a := &API{R: r, s: s, l: logger}
	r.HandleFunc("/hash/{hash}", a.GetAlertByHash)
	r.HandleFunc("/alert/{name}", a.GetRelated)
	r.HandleFunc("/lastupdated", a.GetLastUpdated)
	return a
}

func (a *API) alertHashWriter(w http.ResponseWriter, r *http.Request, h string) {
	if data, ok := a.s.Alerts[h]; ok {
		b := []byte{}
		err := json.Unmarshal(b, data)
		if err != nil {
			a.l.Error("failed to unmarshal data", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(b))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// GetAlertByHash returns the alert data according to the specified hash.
func (a *API) GetAlertByHash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a.alertHashWriter(w, r, vars["hash"])
}

// GetRelated gets the list of related alerts according to the name and label set.
func (a *API) GetRelated(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	alert := alerts.NewAlert()
	for key, vals := range params {
		alert.Labels[key] = vals[0]
	}
	h := alert.Hash()
	a.alertHashWriter(w, r, h)
}

// GetLastUpdated gets the time when the Cache was last updated.
func (a *API) GetLastUpdated(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v", a.s.GetLastUpdated().Unix())
}
