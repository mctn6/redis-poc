package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[:6] != "/user/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	idStr := r.URL.Path[6:]
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := GetUserFromCacheOrDB(id)
	if err != nil {
		log.Printf("User not found: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("❌ Failed to load config:", err)
	}

	InitDB(config)
	InitRedis(config)

	mux := http.NewServeMux()
	mux.HandleFunc("/user/", getUserHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:         ":" + config.ServerPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Printf("🚀 Server running on :%s", config.ServerPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("❌ Server error:", err)
	}
}
