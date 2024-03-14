package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
)

// User представляет пользователя.
type User struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Age     string   `json:"age"`
	Friends []string `json:"friends"`
}

var (
	mu      sync.Mutex
	lastID  int
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

	// Обработчик создания пользователя.
	r.Post("/create", createUserHandler)
	r.Get("/create", getAllUsers)

	// Запуск сервера на порту 8080.
	port := 8080
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Server started on %s (Replica %d)", addr, replica)
	log.Fatal(http.ListenAndServe(addr, r))
}

// createUserHandler обработчик создания пользователя.
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Сохранение пользователя в базу данных.
	userID, err := saveUser(newUser)
	if err != nil {
		log.Printf("Error saving user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Отправка ID пользователя и статуса 201.
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": userID})

	// Логирование создания пользователя.
	log.Printf("User created with ID: %d", userID)
}

// getAllUsers обработчик для GET-запросов (возвращает всех пользователей).
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := fetchAllUsers()
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Отправка данных всех пользователей в формате JSON с отдельной строкой для каждого пользователя.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	for _, user := range users {
		json.NewEncoder(w).Encode(user)
		fmt.Fprintln(w) // Добавляем новую строку после каждого пользователя.
	}

	// Логирование получения всех пользователей.
	log.Println("All users retrieved")
}

// fetchAllUsers возвращает всех пользователей из базы данных.
func fetchAllUsers() ([]User, error) {
	rows, err := db.Query("SELECT id, name, age, friends FROM users")
	if err != nil {
		return nil, fmt.Errorf("error querying users: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var friendsStr string
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &friendsStr); err != nil {
			return nil, fmt.Errorf("error scanning user: %v", err)
		}

		// Преобразование строки с друзьями в массив строк.
		user.Friends = strings.Split(friendsStr, ",")

		users = append(users, user)
	}

	return users, nil
}

// saveUser сохраняет пользователя в базе данных и возвращает ID пользователя.
func saveUser(newUser User) (int, error) {
	mu.Lock()
	defer mu.Unlock()

	// Вставка данных пользователя в базу данных.
	result, err := db.Exec("INSERT INTO users (name, age, friends) VALUES (?, ?, ?)",
		newUser.Name, newUser.Age, strings.Join(newUser.Friends, ","))
	if err != nil {
		return 0, fmt.Errorf("error inserting user: %v", err)
	}

	// Получение ID пользователя из результата вставки.
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %v", err)
	}

	return int(userID), nil
}
