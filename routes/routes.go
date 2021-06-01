package routes

import (
	"fmt"
	"github/sing3demons/go_mux_api/config"
	"github/sing3demons/go_mux_api/controllers"
	"github/sing3demons/go_mux_api/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func Serve(r *mux.Router) {
	db := config.GetDB()
	v1 := "/api/v1"
	authenticate := middleware.AuthMiddleware
	authorize := middleware.Authorize

	authGroup := fmt.Sprintf(v1 + "/auth")
	authController := controllers.Auth{DB: db}
	{
		r.HandleFunc(authGroup+"/sign-up", authController.SignUp).Methods(http.MethodPost)
		r.HandleFunc(authGroup+"/sign-in", middleware.SignIn).Methods(http.MethodPost)
	}
	secureAuth := r.PathPrefix(authGroup).Subrouter()
	secureAuth.Use(authenticate)

	{
		secureAuth.HandleFunc("/profile", authController.GetProfile).Methods("GET")
		secureAuth.HandleFunc("/profile", authController.UpdateImageProfile).Methods(http.MethodPatch)
		secureAuth.HandleFunc("/profile", authController.UpdateProfile).Methods("PUT")
	}

	//products
	productsController := controllers.Product{DB: db}
	productsGroup := fmt.Sprintf(v1 + "/products")
	secureProduct := r.PathPrefix(productsGroup).Subrouter()
	secureProduct.Use(authenticate)
	secureProduct.Use(authorize)
	{
		r.HandleFunc(productsGroup, productsController.FindAll).Methods(http.MethodGet)
		r.HandleFunc(productsGroup+"/{id}", productsController.FindOne).Methods(http.MethodGet)
		secureProduct.HandleFunc("/{id}", productsController.Update).Methods(http.MethodPut)
		secureProduct.HandleFunc("/{id}", productsController.Delete).Methods(http.MethodDelete)
		secureProduct.HandleFunc("", productsController.Create).Methods(http.MethodPost)
	}

	usersController := controllers.Users{DB: db}
	usersGroup := fmt.Sprintf(v1 + "/users")

	secureUsers := r.PathPrefix(usersGroup).Subrouter()
	secureUsers.Use(authenticate)
	secureUsers.Use(authorize)
	{
		secureUsers.HandleFunc("", usersController.FindAll).Methods(http.MethodGet)
		secureUsers.HandleFunc("", usersController.Create).Methods(http.MethodPost)
		secureUsers.HandleFunc("/{id}", usersController.FindOne).Methods(http.MethodGet)
		secureUsers.HandleFunc("/{id}", usersController.Update).Methods(http.MethodPut)
		secureUsers.HandleFunc("/{id}", usersController.Delete).Methods(http.MethodDelete)
		secureUsers.HandleFunc("/{id}/promote", usersController.Promote).Methods(http.MethodPatch)
		secureUsers.HandleFunc("/{id}/demote", usersController.Demote).Methods(http.MethodPatch)
	}

}
