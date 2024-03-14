package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
)

// User структура представляет пользователя.
type User struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Age     string   `json:"age"`
	Friends []string `json:"friends"`
}

var (
	mu      sync.Mutex
	db      *sql.DB
	replica int
)

func init() {
	// Инициализация базы данных (SQLite в данном случае).
	var err error
	db, err = sql.Open("sqlite3", "users.db")
	if err != nil {
		log.Fatal(err)
	}

	// Создание таблицы, если она не существует.
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		name TEXT,
		age TEXT,
		friends TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Настройка маршрутизатора Chi.
	r := chi.NewRouter()

	// Обработчик удаления пользователя.
	r.Delete("/user", deleteUser)

	// Запуск сервера на порту 8081.
	port := 8081
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Server started on %s (Replica %d)", addr, replica)
	log.Fatal(http.ListenAndServe(addr, r))
}

// deleteUser обработчик удаления пользователя.
func deleteUser(w http.ResponseWriter, r *http.Request) {
	var targetID struct {
		TargetID int `json:"target_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&targetID); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Удаление пользователя из базы данных.
	userName, err := deleteUserByID(targetID.TargetID)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Отправка успешного ответа с именем удаленного пользователя.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": fmt.Sprintf("User %s deleted", userName)})

	// Логирование удаления пользователя.
	log.Printf("User with ID %d deleted", targetID.TargetID)
}

// deleteUserByID удаляет пользователя из базы данных по его ID и возвращает его имя.
func deleteUserByID(userID int) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	var userName string
	err := db.QueryRow("SELECT name FROM users WHERE id = ?", userID).Scan(&userName)
	if err != nil {
		return "", fmt.Errorf("error querying user: %v", err)
	}

	_, err = db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return "", fmt.Errorf("error deleting user: %v", err)
	}

	return userName, nil
}
