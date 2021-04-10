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
		r.HandleFunc(authGroup+"/sign-in", controllers.SignIn)
	}

	//products
	productsController := controllers.Product{DB: db}
	productsGroup := fmt.Sprintf(v1 + "/products")

	secureProduct := r.PathPrefix(productsGroup).Subrouter()
	secureProduct.Use(controllers.Middleware)

	{
		r.HandleFunc(productsGroup, productsController.FindAll).Methods(http.MethodGet)
		r.HandleFunc(productsGroup+"/{id}", productsController.FindOne).Methods(http.MethodGet)
		secureProduct.HandleFunc("/{id}", productsController.Update).Methods(http.MethodPut)
		secureProduct.HandleFunc("/{id}", productsController.Delete).Methods(http.MethodDelete)
		secureProduct.HandleFunc("", productsController.Create).Methods(http.MethodPost)
	}

}
