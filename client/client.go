package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Transaction struct {
	ID     int     `json:"id"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
}

func add(transaction Transaction) (*Transaction, error) {
	//Ссылаемся на адрес сервера
	url := "http://localhost:8080/addTransaction"
	//Сереализуем сообщ
	data, err := json.Marshal(transaction)
	if err != nil {
		return nil, err
	}
	//Отправлем метод post
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	//Закрываем тело запроса
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Ошибка добавления транзакции:%s", resp.Status)
	}
	var addedTransaction Transaction
	err = json.NewDecoder(resp.Body).Decode(&addedTransaction)
	if err != nil {
		return nil, err
	}
	return &addedTransaction, nil

}

func get() ([]Transaction, error) {
	url := "http://localhost:8080/getTransaction"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ошибка получения транзакций")
	}
	var getransaction []Transaction
	err = json.NewDecoder(resp.Body).Decode(&getransaction)
	if err != nil {
		return nil, err
	}
	return getransaction, nil
}

func main() {
	NewTransaction := Transaction{Amount: 100.0, Type: "income"}
	add, err := add(NewTransaction)
	if err != nil {
		fmt.Println("Ошибка", err)
	} else {
		fmt.Print("Добавлена транзакция:", add)
	}
	transactions, err := get()
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Список транзакций:", transactions)
	}

}
