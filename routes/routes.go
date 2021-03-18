package routes

import (
	"app/config"
	"app/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func Serve(r *mux.Router) {
	db := config.GetDB()
	v1 := "/api/v1"

	//products
	productsController := controllers.Product{DB: db}
	productsGroup := v1 + "/products"
	r.HandleFunc(productsGroup, productsController.FindAll).Methods(http.MethodGet)
	r.HandleFunc(productsGroup+"/{id}", productsController.FindOne).Methods(http.MethodGet)
	r.HandleFunc(productsGroup, productsController.Create).Methods(http.MethodPost)

}
