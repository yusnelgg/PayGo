package services

import (
	"errors"
	"time"

	"paygo/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("el usuario ya existe")
	ErrUserNotFound      = errors.New("usuario no encontrado")
	ErrInvalidPassword   = errors.New("password incorrecto")
	ErrInvalidToken      = errors.New("token inválido")
)

var users []models.User
var idCounter int = 1

var jwtSecret = []byte("paygo_secret_key_change_in_production")

func InitUserService() {
	idCounter = 1
	users = []models.User{}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func RegisterUser(req models.RegisterRequest) (models.User, error) {
	if BuscarUsuarioPorEmail(req.Email) != nil {
		return models.User{}, ErrUserAlreadyExists
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		ID:           idCounter,
		Username:     req.Username,
		Email:        req.Email,
		Password:     req.Password,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	idCounter++
	users = append(users, user)

	return user, nil
}

func BuscarUsuarioPorEmail(email string) *models.User {
	for i := range users {
		if users[i].Email == email {
			return &users[i]
		}
	}
	return nil
}

func BuscarUsuarioPorID(id int) *models.User {
	for i := range users {
		if users[i].ID == id {
			return &users[i]
		}
	}
	return nil
}

func LoginUser(req models.LoginRequest) (models.LoginResponse, error) {
	user := BuscarUsuarioPorEmail(req.Email)
	if user == nil {
		return models.LoginResponse{}, ErrUserNotFound
	}

	if !CheckPassword(req.Password, user.PasswordHash) {
		return models.LoginResponse{}, ErrInvalidPassword
	}

	token, err := GenerateToken(user.ID, user.Email)
	if err != nil {
		return models.LoginResponse{}, err
	}

	userSinPassword := *user
	userSinPassword.Password = ""
	userSinPassword.PasswordHash = ""

	return models.LoginResponse{
		Token: token,
		User:  userSinPassword,
	}, nil
}

func GetUserProfile(userID int) (models.User, error) {
	user := BuscarUsuarioPorID(userID)
	if user == nil {
		return models.User{}, ErrUserNotFound
	}

	userSinPassword := *user
	userSinPassword.Password = ""
	userSinPassword.PasswordHash = ""

	return userSinPassword, nil
}

func ListarUsuarios() []models.User {
	result := make([]models.User, len(users))
	for i, u := range users {
		u.Password = ""
		u.PasswordHash = ""
		result[i] = u
	}
	return result
}
