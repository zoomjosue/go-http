package main

import (
	"fmt"
	"log"
	"net/http"

	"go-http/handlers"
)

func main() {
	puerto := "24918"

	http.HandleFunc("/api/dragones", handlers.RutaDragones)
	http.HandleFunc("/api/dragones/", handlers.RutaDragones)

	fmt.Printf("API HTTYD corriendo en http://localhost:%s\n", puerto)
	fmt.Println("Endpoints disponibles:")
	fmt.Println("  GET    /api/dragones")
	fmt.Println("  GET    /api/dragones?id=1&especie=X&jinete=X&es_alfa=true&ubicacion=X")
	fmt.Println("  GET    /api/dragones/{id}")
	fmt.Println("  POST   /api/dragones")
	fmt.Println("  PUT    /api/dragones/{id}")
	fmt.Println("  PATCH  /api/dragones/{id}")
	fmt.Println("  DELETE /api/dragones/{id}")

	log.Fatal(http.ListenAndServe(":"+puerto, nil))
}
