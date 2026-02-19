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
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer body", http.StatusBadRequest)
		return
	}
	var user models.User
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	usuarioRegistrado := services.RegistrarUsuario(user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarioRegistrado)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer body", http.StatusBadRequest)
		return
	}

	var cred struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(body, &cred); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	user := services.BuscarUsuario(cred.Email)

	// validate result
	if user == nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	if user.Password == cred.Password {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	} else {
		http.Error(w, "Password incorrecto", http.StatusUnauthorized)
	}
}
func listarUsuarios(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		usuarios := services.ListarUsuarios()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(usuarios)
	case http.MethodPost:
		// could forward to registrar or similar logic
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "Crear usuario no implementado")
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
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
