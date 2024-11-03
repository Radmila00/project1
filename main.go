package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v4"
)

var db *pgx.Conn

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

var mutex sync.Mutex //для безопасного доступа нескольких горутин

func main() {
	var err error
	db, err = pgx.Connect(context.Background(), "postgres://radmila:postgres@localhost:5432/trekerTransaction")
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных", err)
	}
	defer db.Close(context.Background())
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/addTransaction", addTransaction)
	mux.HandleFunc("/getTransaction", getTransaction)
	fmt.Println("Сервер запущен на порту 8080")
	er := http.ListenAndServe(":8080", mux)
	if er != nil {
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

	var id int
	mutex.Lock()
	err := db.QueryRow(context.Background(), "INSERT INTO Transaction (amount,type) VALUES ($1,$2) RETURNING id", t.Amount, t.Type).Scan(&id)

	if err != nil {
		log.Println("Ошибка при добавлении транзакции:", err)
		http.Error(w, "Ошибка добавления транзакции", http.StatusInternalServerError)
		return
	}
	defer mutex.Unlock()
	t.ID = id

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func getTransaction(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	mutex.Lock()
	rows, err := db.Query(context.Background(), "SELECT id, amount,type FROM Transaction")

	if err != nil {
		http.Error(w, "Ошибка получения транзакций", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	defer mutex.Unlock()
	var transaction []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.Amount, &t.Type); err != nil {
			http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
			return
		}
		transaction = append(transaction, t)
	}
	json.NewEncoder(w).Encode(transaction)

}
