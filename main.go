package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Transaction struct {
	ID     int     `json:"id"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	transaction []Transaction
	nextID      = 1
	mutex       sync.Mutex //для безопасного доступа нескольких горутин
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/addTransaction", addTransaction)
	mux.HandleFunc("/getTransaction", getTransaction)
	fmt.Println("Сервер запущен на порту 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Это главная страница"))
}

func addTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var t Transaction
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}
	mutex.Lock()
	t.ID = nextID
	nextID++
	transaction = append(transaction, t)
	mutex.Unlock()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func getTransaction(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}
