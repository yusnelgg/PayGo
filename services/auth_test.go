package services

import (
	"paygo/models"
	"reflect"
	"testing"
)

func reset() {
	// reset global state between tests
	users = nil
	idCounter = 1
}

func TestRegistrarYBuscarUsuario(t *testing.T) {
	reset()

	u := models.User{Username: "test", Email: "a@b.com", Password: "pass"}
	registered := RegistrarUsuario(u)
	if registered.ID != 1 {
		t.Errorf("expected ID 1, got %d", registered.ID)
	}

	found := BuscarUsuario("a@b.com")
	if found == nil {
		t.Fatal("expected to find user, got nil")
	}
	if found.Email != u.Email {
		t.Errorf("expected email %s, got %s", u.Email, found.Email)
	}
}

func TestBuscarUsuario_NoExist(t *testing.T) {
	reset()
	if user := BuscarUsuario("no@existe"); user != nil {
		t.Errorf("expected nil, got %+v", user)
	}
}

func TestListarUsuarios(t *testing.T) {
	reset()
	u1 := RegistrarUsuario(models.User{Username: "u1", Email: "1@a.com", Password: "p1"})
	u2 := RegistrarUsuario(models.User{Username: "u2", Email: "2@a.com", Password: "p2"})

	list := ListarUsuarios()
	if len(list) != 2 {
		t.Fatalf("expected 2 users, got %d", len(list))
	}
	// compare by ignoring order
	want := []models.User{u1, u2}
	if !reflect.DeepEqual(list, want) {
		t.Errorf("list mismatch\nwant %+v\ngot  %+v", want, list)
	}
}
