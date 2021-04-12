package main

import (
	"app/config"
	"app/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDB()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World"))
	}).Methods(http.MethodGet)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Choose the folder to serve
	staticDir := "uploads/"

	uploadDir := [...]string{"products", "users"}
	for _, path := range uploadDir {
		os.MkdirAll(staticDir+path, 0755)
	}

	// Create the route
	r.PathPrefix("/" + staticDir).Handler(http.StripPrefix("/"+staticDir, http.FileServer(http.Dir("./"+staticDir))))
	routes.Serve(r)

	port := fmt.Sprintf(":" + os.Getenv("PORT"))
	fmt.Printf("Running on port %s \n", port)

	http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(r))
}
