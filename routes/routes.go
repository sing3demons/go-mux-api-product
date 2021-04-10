package routes

import (
	"app/config"
	"app/controllers"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Serve(r *mux.Router) {
	db := config.GetDB()
	v1 := "/api/v1"

	authGroup := fmt.Sprintf(v1 + "/auth")
	authController := controllers.Auth{DB: db}
	{
		r.HandleFunc(authGroup+"/sign-up", authController.SignUp)
	}

	//products
	productsController := controllers.Product{DB: db}
	productsGroup := fmt.Sprintf(v1 + "/products")
	{
		r.HandleFunc(productsGroup, productsController.FindAll).Methods(http.MethodGet)
		r.HandleFunc(productsGroup+"/{id}", productsController.FindOne).Methods(http.MethodGet)
		r.HandleFunc(productsGroup+"/{id}", productsController.Update).Methods(http.MethodPut)
		r.HandleFunc(productsGroup+"/{id}", productsController.Delete).Methods(http.MethodDelete)
		r.HandleFunc(productsGroup, productsController.Create).Methods(http.MethodPost)
	}

}
