package services

import (
	"fmt"
	"paygo/models"
)

var users []models.User
var idCounter int = 1

func RegistrarUsuario(user models.User) models.User {
	user.ID = idCounter
	idCounter++
	users = append(users, user)
	fmt.Printf("Usuario registrado: %s\n", user.Username)
	return user
}
