package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"quiz-go/src/infrastructure/database/mongodb"
)

type Question struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Question QuestionDetails    `json:"question"`
}

type QuestionDetails struct {
	Description   string   `json:"description"`
	Explanation   string   `json:"explanation"`
	Difficulty    string   `json:"difficulty"`
	Categories    []string `json:"categories"`
	AllowMultiple bool     `json:"allow_multiple"`
	Options       []Option `json:"options"`
}

type Option struct {
	OptionText string `json:"option_text"`
	IsCorrect  bool   `json:"is_correct"`
}

type Category struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name string             `json:"name"`
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleError(w http.ResponseWriter, err error, message string, statusCode int) {
	log.Printf("%s: %v", message, err)
	http.Error(w, message, statusCode)
}

func convertID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Erro ao codificar resposta em JSON: %v", err)
		http.Error(w, "Erro ao retornar resposta", http.StatusInternalServerError)
	}
}

func getQuestionByID(id string) (Question, error) {
	var question Question
	objectID, err := convertID(id)
	if err != nil {
		return question, err
	}

	filter := bson.M{"_id": objectID}
	err = mongodb.QuestionCollection.FindOne(context.TODO(), filter).Decode(&question)
	return question, err
}

func getQuestion(w http.ResponseWriter, r *http.Request) {
	questionID := mux.Vars(r)["id"]

	question, err := getQuestionByID(questionID)
	if err != nil {
		handleError(w, err, "Erro ao buscar a questão", http.StatusNotFound)
		return
	}

	log.Printf("Get Questão com ID %s", questionID)
	sendJSONResponse(w, question, http.StatusOK)
}

func deleteQuestion(w http.ResponseWriter, r *http.Request) {
	questionID := mux.Vars(r)["id"]

	question, err := getQuestionByID(questionID)
	if err != nil {
		handleError(w, err, "Erro ao buscar a questão", http.StatusNotFound)
		return
	}

	// Remove the question
	filter := bson.M{"_id": question.ID}
	deleteResult, err := mongodb.QuestionCollection.DeleteOne(context.TODO(), filter)
	if err != nil || deleteResult.DeletedCount == 0 {
		handleError(w, err, "Erro ao remover a questão", http.StatusInternalServerError)
		return
	}

	log.Printf("Get Questão com ID %s", questionID)
	w.WriteHeader(http.StatusNoContent)
}

func updateQuestion(w http.ResponseWriter, r *http.Request) {
	questionID := mux.Vars(r)["id"]

	var updatedDetails QuestionDetails
	if err := json.NewDecoder(r.Body).Decode(&updatedDetails); err != nil {
		handleError(w, err, "Erro no corpo da requisição", http.StatusBadRequest)
		return
	}

	objectID, err := convertID(questionID)
	if err != nil {
		handleError(w, err, "ID inválido", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"question": updatedDetails}}

	result := mongodb.QuestionCollection.FindOneAndUpdate(context.TODO(), filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		handleError(w, result.Err(), "Erro ao atualizar questão", http.StatusInternalServerError)
		return
	}

	var updatedQuestion Question
	if err := result.Decode(&updatedQuestion); err != nil {
		handleError(w, err, "Erro ao recuperar questão atualizada", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, updatedQuestion, http.StatusOK)
}

func createQuestion(w http.ResponseWriter, r *http.Request) {
	var questionDetails QuestionDetails
	if err := json.NewDecoder(r.Body).Decode(&questionDetails); err != nil {
		handleError(w, err, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	newQuestion := Question{
		ID:       primitive.NewObjectID(),
		Question: questionDetails,
	}

	_, err := mongodb.QuestionCollection.InsertOne(context.TODO(), newQuestion)
	if err != nil {
		handleError(w, err, "Erro ao inserir questão no MongoDB", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, newQuestion, http.StatusCreated)
}

func getQuestions(w http.ResponseWriter, r *http.Request) {
	var questions []Question
	cursor, err := mongodb.QuestionCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		handleError(w, err, "Erro ao buscar perguntas", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var question Question
		if err := cursor.Decode(&question); err != nil {
			handleError(w, err, "Erro ao decodificar pergunta", http.StatusInternalServerError)
			return
		}
		questions = append(questions, question)
	}

	sendJSONResponse(w, questions, http.StatusOK)
}

// CRUD Operations for Categories
func createCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		handleError(w, err, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	category.ID = primitive.NewObjectID()
	_, err := mongodb.CategoryCollection.InsertOne(context.TODO(), category)
	if err != nil {
		handleError(w, err, "Erro ao inserir categoria", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, category, http.StatusCreated)
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	var categories []Category
	cursor, err := mongodb.CategoryCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		handleError(w, err, "Erro ao buscar categorias", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var category Category
		if err := cursor.Decode(&category); err != nil {
			handleError(w, err, "Erro ao decodificar categoria", http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	sendJSONResponse(w, categories, http.StatusOK)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := mux.Vars(r)["id"]

	var updatedCategory Category
	if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
		handleError(w, err, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	objectID, err := convertID(categoryID)
	if err != nil {
		handleError(w, err, "ID inválido", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"name": updatedCategory.Name}}

	result := mongodb.CategoryCollection.FindOneAndUpdate(context.TODO(), filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		handleError(w, result.Err(), "Erro ao atualizar categoria", http.StatusInternalServerError)
		return
	}

	var updatedCategoryResponse Category
	if err := result.Decode(&updatedCategoryResponse); err != nil {
		handleError(w, err, "Erro ao recuperar categoria atualizada", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, updatedCategoryResponse, http.StatusOK)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := mux.Vars(r)["id"]

	objectID, err := convertID(categoryID)
	if err != nil {
		handleError(w, err, "ID inválido", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objectID}
	deleteResult, err := mongodb.CategoryCollection.DeleteOne(context.TODO(), filter)
	if err != nil || deleteResult.DeletedCount == 0 {
		handleError(w, err, "Erro ao remover categoria", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	mongodb.Connect()

	r := mux.NewRouter()
	log.Println("Configurando rotas")
	r.HandleFunc("/questions/{id}", deleteQuestion).Methods("DELETE")
	r.HandleFunc("/questions", createQuestion).Methods("POST")
	r.HandleFunc("/questions", getQuestions).Methods("GET")
	r.HandleFunc("/questions/{id}", updateQuestion).Methods("PUT")
	r.HandleFunc("/questions/{id}", getQuestion).Methods("GET")

	// Category routes
	r.HandleFunc("/categories", createCategory).Methods("POST")
	r.HandleFunc("/categories", getCategories).Methods("GET")
	r.HandleFunc("/categories/{id}", updateCategory).Methods("PUT")
	r.HandleFunc("/categories/{id}", deleteCategory).Methods("DELETE")

	log.Println("Iniciando servidor na porta 8000")
	http.Handle("/", enableCORS(r))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
