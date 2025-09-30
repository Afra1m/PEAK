package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"PEAK/internal/db"
	"PEAK/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(func() []byte {
	if v := os.Getenv("JWT_SECRET"); v != "" {
		return []byte(v)
	}
	return []byte("change-this-secret")
}())

type credsReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var cr credsReq
	if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	// hash
	hash, err := bcrypt.GenerateFromPassword([]byte(cr.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	// insert
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	q := `INSERT INTO users (email, password_hash, full_name, role_id) VALUES ($1,$2,$3,$4) RETURNING id, created_at`
	var id string
	var createdAt time.Time
	roleId := 1 // default engineer (or choose)
	err = db.Pool.QueryRow(ctx, q, cr.Email, string(hash), cr.Name, roleId).Scan(&id, &createdAt)
	if err != nil {
		http.Error(w, "could not create", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"id": id, "created_at": createdAt})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var cr credsReq
	if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var user models.User
	q := `SELECT id, password_hash, role_id FROM users WHERE email=$1`
	row := db.Pool.QueryRow(ctx, q, cr.Email)
	if err := row.Scan(&user.ID, &user.PasswordHash, &user.RoleID); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(cr.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	// create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.RoleID,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})
	tokStr, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": tokStr})
}
