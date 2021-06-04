package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int) func(v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	w.WriteHeader(statusCode)
	return func(v interface{}) error {
		return json.NewEncoder(w).Encode(v)
	}
}

func BodyParser(r *http.Request) []byte {
	body, _ := ioutil.ReadAll(r.Body)
	return body
}

func ToJson(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-type", "application/json; charset=UTF8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	CheckError(err)
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
