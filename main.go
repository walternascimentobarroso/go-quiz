package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"quiz-go/src/handlers"
	"quiz-go/src/infrastructure/database/mongodb"
	"quiz-go/src/middlewares"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	mongodb.Connect()

	r := mux.NewRouter()
	log.Println("Configurando rotas")
	r.HandleFunc("/questions/{id}", handlers.DeleteQuestion).Methods("DELETE")
	r.HandleFunc("/questions", handlers.CreateQuestion).Methods("POST")
	r.HandleFunc("/questions", handlers.GetQuestions).Methods("GET")
	r.HandleFunc("/questions/{id}", handlers.UpdateQuestion).Methods("PUT")
	r.HandleFunc("/questions/{id}", handlers.GetQuestion).Methods("GET")

	// Category routes
	r.HandleFunc("/categories", handlers.CreateCategory).Methods("POST")
	r.HandleFunc("/categories", handlers.GetCategories).Methods("GET")
	r.HandleFunc("/categories/{id}", handlers.UpdateCategory).Methods("PUT")
	r.HandleFunc("/categories/{id}", handlers.DeleteCategory).Methods("DELETE")

	log.Println("Iniciando servidor na porta 8000")
	http.Handle("/", middlewares.EnableCORS(r))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
