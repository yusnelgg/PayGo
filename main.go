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

func login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer body", 400)
		return
	}

	var cred struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(body, &cred); err != nil {
		http.Error(w, "JSON inválido", 400)
		return
	}

	user := services.BuscarUsuario(cred.Email)

	if user.ID == 0 {
		http.Error(w, "Usuario no encontrado", 404)
		return
	}

	if user.Password == cred.Password {
		json.NewEncoder(w).Encode(user)
	} else {
		http.Error(w, "Password incorrecto", 401)
	}
}
func listarUsuarios(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		usuarios := services.ListarUsuarios()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(usuarios)
	case "POST":
		fmt.Fprintf(w, "Crear usuario")

	}
}

func main() {
	fmt.Println("PayGo iniciando en http://localhost:8080")

	http.HandleFunc("/", inicio)
	http.HandleFunc("/login", login)
	http.HandleFunc("/registrar", registrar)
	http.HandleFunc("/users", listarUsuarios)

	http.ListenAndServe(":8080", nil)
}
