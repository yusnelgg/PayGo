package services

import (
	"fmt"
	"paygo/models"
)

// package-level state; tests can reset via UsersReset
var users []models.User
var idCounter int = 1

func RegistrarUsuario(user models.User) models.User {
	user.ID = idCounter
	idCounter++
	users = append(users, user)
	fmt.Printf("Usuario registrado: %s\n", user.Username)
	return user
}

func BuscarUsuario(email string) *models.User {
	// iterate by index so we can return the address of the element
	for i := range users {
		if users[i].Email == email {
			return &users[i]
		}
	}
	return nil
}

func ListarUsuarios() []models.User {
	return users
}

// UsersReset clears in-memory storage; primarily for tests.
func UsersReset() {
	users = nil
	idCounter = 1
}
