package main

import (
	"log"
	"net/http"

	"daifugo-backend/internal/server"
)

func main() {
	gameManager := server.NewGameManager()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWs(gameManager, w, r)
	})

	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)

	log.Println("Starting Server on :5001")
	if err := http.ListenAndServe(":5001", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
