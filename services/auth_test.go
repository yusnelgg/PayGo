package services

import (
	"paygo/models"
	"testing"
)

func TestRegistrarYBuscarUsuario(t *testing.T) {
	InitUserService()

	req := models.RegisterRequest{
		Username: "test",
		Email:    "a@b.com",
		Password: "password123",
	}
	registered, err := RegisterUser(req)
	if err != nil {
		t.Fatalf("error registering: %v", err)
	}
	if registered.ID != 1 {
		t.Errorf("expected ID 1, got %d", registered.ID)
	}

	found := BuscarUsuarioPorEmail("a@b.com")
	if found == nil {
		t.Fatal("expected to find user, got nil")
	}
	if found.Email != req.Email {
		t.Errorf("expected email %s, got %s", req.Email, found.Email)
	}
}

func TestBuscarUsuario_NoExist(t *testing.T) {
	InitUserService()
	if user := BuscarUsuarioPorEmail("no@existe"); user != nil {
		t.Errorf("expected nil, got %+v", user)
	}
}

func TestListarUsuarios(t *testing.T) {
	InitUserService()
	RegisterUser(models.RegisterRequest{Username: "u1", Email: "1@a.com", Password: "p1"})
	RegisterUser(models.RegisterRequest{Username: "u2", Email: "2@a.com", Password: "p2"})

	list := ListarUsuarios()
	if len(list) != 2 {
		t.Fatalf("expected 2 users, got %d", len(list))
	}
}

func TestLoginUser(t *testing.T) {
	InitUserService()

	RegisterUser(models.RegisterRequest{
		Username: "test",
		Email:    "test@b.com",
		Password: "password123",
	})

	_, err := LoginUser(models.LoginRequest{
		Email:    "test@b.com",
		Password: "wrong",
	})
	if err != ErrInvalidPassword {
		t.Errorf("expected ErrInvalidPassword, got %v", err)
	}

	resp, err := LoginUser(models.LoginRequest{
		Email:    "test@b.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("error logging in: %v", err)
	}
	if resp.Token == "" {
		t.Error("expected token, got empty")
	}
}
