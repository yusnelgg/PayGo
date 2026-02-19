package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"paygo/models"
	"paygo/services"
)

func inicio(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola, has accedido a %s", r.URL.Path)
}

func registrar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error al leer body", 400)
			return
		}
		var user models.User
		if err := json.Unmarshal(body, &user); err != nil {
			http.Error(w, "JSON inválido", 400)
			return
		}
		usuarioRegistrado := services.RegistrarUsuario(user)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(usuarioRegistrado)

	}
}

func pagos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Ver pagos")
	case "POST":
		fmt.Fprintf(w, "Crear pago")
	}
}

func main() {
	fmt.Println("PayGo iniciando en http://localhost:8080")

	http.HandleFunc("/", inicio)
	http.HandleFunc("/registrar", registrar)
	http.HandleFunc("/pagos", pagos)

	http.ListenAndServe(":8080", nil)
}
