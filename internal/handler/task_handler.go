package handler

import (
	"database/sql"
	"encoding/json"
	"microservices/internal/db"
	"net/http"
)

type TaskHandler struct {
	Q *db.Queries
}

func NewTaskHandler(q *db.Queries) *TaskHandler {
	return  &TaskHandler{q}
}

///------------ CURD ----------------

//  ========= Create ========

func(h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title		    string    `json:"title"`
		Description string    `json:"description"`
		StatusID    int32     `json:"status_id"`
		UserID      int32     `json:"user_id"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request", 400)
		return
	}

	// Validation
	

	if req.Title
  _, err = h.Q.CreateTask(r.Context(), db.CreateTaskParams{
		Title:req.Title,
		Description: sql.NullString{
        String: req.Description,
				Valid: req.Description != "",
		},
		Status: sql.NullString{
			String: "pending",
			Valid: true,
		},
	})
	if err != nil {
		http.Error(w, "Failed to create", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

