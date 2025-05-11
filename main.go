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
func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func main() {
	http.HandleFunc("GET /todos", handleGetTodos)
	http.HandleFunc("POST /todos", handleCreateTodo)
	http.HandleFunc("PUT /todos", handleCompleteFunc)
	http.ListenAndServe(":"+getEnv("PORT", "8080"), nil)
}
