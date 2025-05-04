package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

type Todo struct {
	ID        int
	Title     string
	Completed bool
}

var todos = []Todo{
}

func handleGetTodos(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(todos)
}

func handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todos = append(todos, newTodo)
	json.NewEncoder(w).Encode(newTodo)
}

func handleCompleteFunc(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if idInt >= len(todos) || idInt < 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	todos[idInt].Completed = true
	json.NewEncoder(w).Encode(todos[idInt])
}

func handleTodos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetTodos(w, r)
	case http.MethodPost:
		handleCreateTodo(w, r)
	case http.MethodPut:
		handleCompleteFunc(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	http.HandleFunc("/todos", handleTodos)
	http.ListenAndServe(":"+port, nil)
}

// testing 
// curl -X GET http://localhost:8080/todos
// curl -X POST http://localhost:8080/todos -H "Content-Type: application/json" -d '{"title": "Buy groceries"}'
// curl -X PUT http://localhost:8080/todos?id=0
