package main

import (
	"app/config"
	"app/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var port string = "8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World"))
	}).Methods(http.MethodGet)

	config.InitDB()

	// Choose the folder to serve
	staticDir := "uploads/"

	uploadDir := [...]string{"products", "users"}
	for _, path := range uploadDir {

		os.MkdirAll(staticDir+path, 0755)
	}

	// Create the route
	r.PathPrefix("/" + staticDir).Handler(http.StripPrefix("/"+staticDir, http.FileServer(http.Dir("./"+staticDir))))
	routes.Serve(r)

	fmt.Printf("Running on port :%s \n", os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}
