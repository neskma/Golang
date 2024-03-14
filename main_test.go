package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestNameOnCreatePage(t *testing.T) {
	// Выполнение GET-запроса к серверу на порту 8080
	res, err := http.Get("http://localhost:8080/create")
	if err != nil {
		t.Fatalf("Error performing GET request: %v", err)
	}
	defer res.Body.Close()

	// Проверка кода состояния ответа
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d but got %d", http.StatusOK, res.StatusCode)
	}

	// Чтение содержимого ответа
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	// Проверка наличия имени в теле ответа
	expectedName := "Timothy"
	if !strings.Contains(string(body), expectedName) {
		t.Errorf("Expected page body to contain %q, but it didn't", expectedName)
	}
}
