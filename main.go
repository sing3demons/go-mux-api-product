package main

import (
	"encoding/json"
	"flag"
	"github/sing3demons/go_mux_api/config"
	"github/sing3demons/go_mux_api/routes"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var dir string

	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	config.InitDB()

	r := mux.NewRouter()
	r.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "rest api golang"})
	}).Methods(http.MethodGet)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Choose the folder to serve
	staticDir := "uploads"

	uploadDir := [...]string{"products", "users"}
	for _, path := range uploadDir {
		os.MkdirAll(staticDir + "/" +path, 0755)
	}

	// Create the route
	r.PathPrefix("/" + staticDir + "/").Handler(http.StripPrefix("/"+staticDir+"/", http.FileServer(http.Dir("./"+staticDir+"/"))))
	routes.Serve(r)

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	srv := &http.Server{
		Handler:      handlers.CORS(originsOk, headersOk, methodsOk)(loggedRouter),
		Addr:         ":" + os.Getenv("PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}



	log.Printf("Running on port : %s  \n", os.Getenv("PORT"))
	log.Fatal(srv.ListenAndServe())
}
