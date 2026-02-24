package main

import (
	"fmt"
	"net/http"

	"paygo/handlers"
	"paygo/services"
)

func main() {
	services.InitUserService()

	fmt.Println("PayGo API iniciando en http://localhost:8080")
	fmt.Println("Endpoints disponibles:")
	fmt.Println("  POST /api/users/register")
	fmt.Println("  POST /api/users/login")
	fmt.Println("  GET  /api/users/profile (requiere token)")
	fmt.Println("  GET  /api/users (requiere token)")
	fmt.Println("  GET  /api/users?id=1 (requiere token)")

	http.HandleFunc("/api/users/register", handlers.Register)
	http.HandleFunc("/api/users/login", handlers.Login)
	http.HandleFunc("/api/users/profile", handlers.GetProfile)
	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("id") != "" {
			handlers.GetUserByID(w, r)
		} else {
			handlers.ListUsers(w, r)
		}
	})

	http.ListenAndServe(":8080", nil)
}
