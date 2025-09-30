package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"PEAK/internal/db"
	"PEAK/internal/models"

	"github.com/gorilla/mux"
)

func CreateProject(w http.ResponseWriter, r *http.Request) {
	var p models.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	q := `INSERT INTO projects (name, description, created_by) VALUES ($1,$2,$3) RETURNING id, created_at`
	var id string
	var createdAt time.Time
	// created_by should come from JWT in real scenario; here accept passed or from middleware
	createdBy := p.CreatedBy
	if createdBy == "" {
		createdBy = ""
	} // adjust
	if err := db.Pool.QueryRow(ctx, q, p.Name, p.Description, createdBy).Scan(&id, &createdAt); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	p.ID = id
	p.CreatedAt = createdAt
	json.NewEncoder(w).Encode(p)
}

func ListProjects(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	rows, err := db.Pool.Query(ctx, `SELECT id, name, description, created_by, created_at FROM projects ORDER BY created_at DESC`)
	if err != nil {
		http.Error(w, "db", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var res []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedBy, &p.CreatedAt); err != nil {
			http.Error(w, "scan", http.StatusInternalServerError)
			return
		}
		res = append(res, p)
	}
	json.NewEncoder(w).Encode(res)
}

func GetProject(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var p models.Project
	row := db.Pool.QueryRow(ctx, `SELECT id, name, description, created_by, created_at FROM projects WHERE id=$1`, id)
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedBy, &p.CreatedAt); err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(p)
}
