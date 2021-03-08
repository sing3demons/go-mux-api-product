package routes

import (
	"app/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func Serve() *mux.Router {
	v1 := "/api/v1"

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World"))
	}).Methods(http.MethodGet)

	//products
	router.HandleFunc(v1+"/products", controllers.FindAll).Methods(http.MethodGet)

	return router
}
