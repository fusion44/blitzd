package api

import (
	"encoding/json"
	"github.com/frennkie/blitzd/internal/metric"
	"net/http"
)

func System() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		js, err := json.Marshal(metric.System)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(js)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}