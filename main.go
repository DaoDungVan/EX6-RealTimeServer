package main

import (
	"fmt"
	"realtime-server/internal/auth"
	"realtime-server/internal/chat"
	"realtime-server/internal/common/config"
	"realtime-server/internal/common/db"
	"realtime-server/internal/friend"
	"realtime-server/internal/user"
	"realtime-server/pkg/websocket"

	"log"
	"net/http"
)

func main() {
	config.LoadEnv()

	database := db.InitMongo()

	userRepo := user.NewRepository(database)
	userHandler := user.NewHandler(userRepo)
	authHandler := auth.NewHandler(userRepo)

	friendRepo := friend.NewRepository(database)
	friendHandler := friend.NewHandler(friendRepo)

	chatRepo := chat.NewRepository(database)
	websocket.SetChatRepo(chatRepo)
	go websocket.WS.Run()

	port := config.GetEnv("PORT", "")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to RealTimeServer v0.0.1")
	})

	http.HandleFunc("/api/register", userHandler.Register)
	http.HandleFunc("/api/login", authHandler.Login)

	http.HandleFunc("/api/friend/request", auth.
		JWTMiddleware(friendHandler.SendRequest))
	http.HandleFunc("/api/friend/accept", auth.
		JWTMiddleware(friendHandler.AcceptRequest))

	http.HandleFunc("/ws", websocket.ServeWS)

	fmt.Println("Server is running at http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
