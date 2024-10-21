package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Question struct {
  ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
  Question    QuestionDetails    `json:"question"`
}

type QuestionDetails struct {
  Description  string   `json:"description"`
  Explanation string   `json:"explanation"`
  Difficulty  string   `json:"difficulty"`
  Categories  []string `json:"categories"`
  AllowMultiple bool     `json:"allow_multiple"`
  Options     []Option  `json:"options"`
}

type Option struct {
  OptionText string `json:"option_text"`
  IsCorrect  bool   `json:"is_correct"`
}

var client *mongo.Client
var questionCollection *mongo.Collection

func main() {
    // Configurações do logger
    log.SetFlags(log.LstdFlags | log.Lshortfile) // Adiciona timestamp e nome do arquivo

    // Conexão com o MongoDB
    clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
    var err error
    client, err = mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
    }
    log.Println("Conectado ao MongoDB com sucesso")

    questionCollection = client.Database("quizdb").Collection("questions")

    r := mux.NewRouter()
    r.HandleFunc("/questions", createQuestion).Methods("POST")
    r.HandleFunc("/questions", getQuestions).Methods("GET")

    log.Println("Iniciando servidor na porta 8000")
    http.Handle("/", r)
    log.Fatal(http.ListenAndServe(":8000", nil))
}

func createQuestion(w http.ResponseWriter, r *http.Request) {
  var requestBody struct {
      Question QuestionDetails `json:"question"`
  }

  // Decodifica o corpo da solicitação
  if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
      log.Printf("Erro ao decodificar JSON: %v", err)
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
  }

  // Log dos dados recebidos
  log.Printf("Recebido: %+v", requestBody)

  // Gera um novo ObjectID para a questão
  newQuestion := Question{
      ID:       primitive.NewObjectID(),
      Question: requestBody.Question, // Preenche os detalhes da questão
  }

  // Inserir a questão no MongoDB
  _, err := questionCollection.InsertOne(context.TODO(), newQuestion)
  if err != nil {
      log.Printf("Erro ao inserir questão no MongoDB: %v", err)
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }

  log.Println("Questão inserida com sucesso")

  // Define o cabeçalho Content-Type como application/json
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)

  // Retorna a questão criada como resposta em formato JSON
  if err := json.NewEncoder(w).Encode(newQuestion); err != nil {
      log.Printf("Erro ao codificar resposta em JSON: %v", err)
  }
}

func getQuestions(w http.ResponseWriter, r *http.Request) {
  var questions []Question

  // Define o cabeçalho Content-Type como application/json
  w.Header().Set("Content-Type", "application/json")

  cursor, err := questionCollection.Find(context.TODO(), bson.M{})
  if err != nil {
      log.Printf("Erro ao buscar perguntas: %v", err)
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }
  defer cursor.Close(context.TODO())

  for cursor.Next(context.TODO()) {
      var question Question
      if err := cursor.Decode(&question); err != nil {
          log.Printf("Erro ao decodificar pergunta: %v", err)
          http.Error(w, err.Error(), http.StatusInternalServerError)
          return
      }
      questions = append(questions, question)
  }

  log.Printf("Total de perguntas recuperadas: %d", len(questions))

  // Define o status HTTP 200 OK
  w.WriteHeader(http.StatusOK)

  // Retorna a lista de perguntas como resposta em formato JSON
  if err := json.NewEncoder(w).Encode(questions); err != nil {
      log.Printf("Erro ao codificar resposta em JSON: %v", err)
  }
}

