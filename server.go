package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type User struct {
	ID      int
	Name    string
	Age     string
	Friends []string
	mu      sync.Mutex
}

var (
	users     = make(map[int]*User)
	lastIndex = 0
	mu        sync.Mutex
)

func main() {
	fmt.Println("Starting the server...")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HTTP сервер\n"))
	})

	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		var newUser User
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		lastIndex++
		newUser.ID = lastIndex
		users[lastIndex] = &newUser

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": newUser.ID})
	})

	r.Delete("/user/{target_id}", func(w http.ResponseWriter, r *http.Request) {
		targetID := chi.URLParam(r, "target_id")
		targetUserID, err := strconv.Atoi(targetID)
		if err != nil {
			http.Error(w, "Invalid target user ID", http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		targetUser, targetExists := users[targetUserID]
		if !targetExists {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		for _, friendName := range targetUser.Friends {
			friend, friendExists := findUserByName(friendName)
			if friendExists {
				friend.mu.Lock()
				friend.Friends = removeFriend(friend.Friends, targetUser.Name)
				friend.mu.Unlock()
			}
		}

		delete(users, targetUserID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": targetUser.Name + " удален"})
	})

	r.Get("/friends/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "user_id")
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		user, exists := users[userIDInt]
		if !exists {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user.Friends)
	})

	r.Put("/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "user_id")
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		user, exists := users[userIDInt]
		if !exists {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		var updateData map[string]string
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		newAge := updateData["new_age"]
		user.Age = newAge

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Возраст пользователя успешно обновлен")
	})

	r.Post("/make_friends", func(w http.ResponseWriter, r *http.Request) {
		var friendshipData struct {
			SourceID string `json:"source_id"`
			TargetID string `json:"target_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&friendshipData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		sourceID, err := strconv.Atoi(friendshipData.SourceID)
		if err != nil {
			http.Error(w, "Invalid source user ID", http.StatusBadRequest)
			return
		}

		targetID, err := strconv.Atoi(friendshipData.TargetID)
		if err != nil {
			http.Error(w, "Invalid target user ID", http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		sourceUser, sourceExists := users[sourceID]
		targetUser, targetExists := users[targetID]

		if !sourceExists || !targetExists {
			http.Error(w, "One or more users not found", http.StatusNotFound)
			return
		}

		sourceUser.mu.Lock()
		targetUser.mu.Lock()

		sourceUser.Friends = append(sourceUser.Friends, targetUser.Name)
		targetUser.Friends = append(targetUser.Friends, sourceUser.Name)

		sourceUser.mu.Unlock()
		targetUser.mu.Unlock()

		w.WriteHeader(http.StatusOK)
		resultMessage := fmt.Sprintf("%s и %s теперь друзья", sourceUser.Name, targetUser.Name)
		fmt.Fprintln(w, resultMessage)
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

// Вспомогательные функции

func findUserByName(name string) (*User, bool) {
	for _, user := range users {
		if user.Name == name {
			return user, true
		}
	}
	return nil, false
}

func removeFriend(friends []string, friendToRemove string) []string {
	var updatedFriends []string
	for _, friend := range friends {
		if friend != friendToRemove {
			updatedFriends = append(updatedFriends, friend)
		}
	}
	return updatedFriends
}
