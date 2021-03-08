package main

import (
	"app/routes"
	"fmt"
	"net/http"
)

var port string = "8080"

func main() {
	r := routes.Serve()

	fmt.Printf("Running on port :%s", port)
	http.ListenAndServe(":"+port, r)
}
