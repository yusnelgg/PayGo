package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"paygo/models"
	"paygo/services"
	"testing"
)

func init() {
	// reset state once before any tests
}

// Helper to reset services package by exporting a reset function

func TestInicioHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	inicio(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestRegistrarHandler(t *testing.T) {
	services.UsersReset()

	user := models.User{Username: "h", Email: "h@x.com", Password: "pwd"}
	b, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/registrar", bytes.NewReader(b))
	w := httptest.NewRecorder()
	registrar(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var got models.User
	json.NewDecoder(resp.Body).Decode(&got)
	if got.ID == 0 {
		t.Error("expected assigned ID")
	}

	// wrong method should be 405
	req2 := httptest.NewRequest("GET", "/registrar", nil)
	w2 := httptest.NewRecorder()
	registrar(w2, req2)
	if w2.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for wrong method, got %d", w2.Result().StatusCode)
	}
}

func TestLoginHandler(t *testing.T) {
	services.UsersReset()
	u := services.RegistrarUsuario(models.User{Username: "l", Email: "l@x.com", Password: "pw"})

	cred := map[string]string{"email": u.Email, "password": "pw"}
	body, _ := json.Marshal(cred)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	login(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login failed, status %d", resp.StatusCode)
	}

	req2 := httptest.NewRequest("GET", "/login", nil)
	w2 := httptest.NewRecorder()
	login(w2, req2)
	if w2.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for wrong method, got %d", w2.Result().StatusCode)
	}
}
