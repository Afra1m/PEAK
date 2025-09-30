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

func CreateDefect(w http.ResponseWriter, r *http.Request) {
	var d models.Defect
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	q := `INSERT INTO defects (project_id, title, description, priority, status, assignee, due_date, created_by)
          VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, created_at, updated_at`
	var id string
	var createdAt, updatedAt time.Time
	if err := db.Pool.QueryRow(ctx, q, d.ProjectID, d.Title, d.Description, d.Priority, d.Status, d.Assignee, d.DueDate, d.CreatedBy).
		Scan(&id, &createdAt, &updatedAt); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	d.ID = id
	d.CreatedAt = createdAt
	d.UpdatedAt = updatedAt
	json.NewEncoder(w).Encode(d)
}

func ListDefects(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	rows, err := db.Pool.Query(ctx, `SELECT id, project_id, title, description, priority, status, assignee, due_date, created_by, created_at, updated_at FROM defects ORDER BY created_at DESC`)
	if err != nil {
		http.Error(w, "db", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var res []models.Defect
	for rows.Next() {
		var d models.Defect
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.Title, &d.Description, &d.Priority, &d.Status, &d.Assignee, &d.DueDate, &d.CreatedBy, &d.CreatedAt, &d.UpdatedAt); err != nil {
			http.Error(w, "scan", http.StatusInternalServerError)
			return
		}
		res = append(res, d)
	}
	json.NewEncoder(w).Encode(res)
}

func UpdateDefectStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	_, err := db.Pool.Exec(ctx, `UPDATE defects SET status=$1, updated_at=now() WHERE id=$2`, body.Status, id)
	if err != nil {
		http.Error(w, "db", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
