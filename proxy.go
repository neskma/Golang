package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	serverOneURL = "http://localhost:8080"
	serverTwoURL = "http://localhost:8081"
)

func main() {
	// Настройка прокси для serverOne.
	proxyOne := httputil.NewSingleHostReverseProxy(mustParseURL(serverOneURL))
	http.HandleFunc("/create", handleProxy(proxyOne))

	// Настройка прокси для serverTwo.
	proxyTwo := httputil.NewSingleHostReverseProxy(mustParseURL(serverTwoURL))
	http.HandleFunc("/user", handleProxy(proxyTwo))

	// Запуск прокси-сервера на порту 9090.
	port := 9090
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Proxy server started on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleProxy(proxy *httputil.ReverseProxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проксирование запроса к соответствующему серверу.
		proxy.ServeHTTP(w, r)
	}
}

func mustParseURL(rawURL string) *url.URL {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("Error parsing URL %s: %v", rawURL, err)
	}
	return parsedURL
}
