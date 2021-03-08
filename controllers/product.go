package controllers

import (
	"encoding/json"
	"net/http"
)

func FindAll(w http.ResponseWriter, r *http.Request) {
	// msg := map[string]string{"message": "Hello, world"}
	jsonify(w)(map[string]string{"message": "Hello, world"})
}

func jsonify(w http.ResponseWriter) func(interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return func(v interface{}) error {
		return json.NewEncoder(w).Encode(v)
	}
}
