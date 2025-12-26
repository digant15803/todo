package main

import (
	"log"
	"net/http"

	httpHandler "todo/internal/http"
	todo "todo/internal/todo"
)

func main() {
	service := todo.NewService()
	handler := httpHandler.NewHandler(service)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
