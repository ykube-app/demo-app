package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

type Task struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at"`
}

type Handler struct {
	db *sql.DB
}

// NewHandler returns an http.Handler that serves the tasks API.
// Routes:
//
//	GET    /api/tasks        — list all tasks
//	POST   /api/tasks        — create a task
//	PATCH  /api/tasks/{id}   — toggle done
//	DELETE /api/tasks/{id}   — delete a task
func NewHandler(db *sql.DB) http.Handler {
	h := &Handler{db: db}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/tasks", h.tasksCollection)
	mux.HandleFunc("/api/tasks/", h.tasksItem)
	return mux
}

func (h *Handler) tasksCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) tasksItem(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	if id == "" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodPatch:
		h.toggleTask(w, r, id)
	case http.MethodDelete:
		h.deleteTask(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.QueryContext(r.Context(),
		`SELECT id, title, done, created_at FROM tasks ORDER BY created_at`)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			http.Error(w, "scan error", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, t)
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Title == "" {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	var t Task
	err := h.db.QueryRowContext(r.Context(),
		`INSERT INTO tasks (title) VALUES ($1) RETURNING id, title, done, created_at`,
		input.Title,
	).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, t)
}

func (h *Handler) toggleTask(w http.ResponseWriter, r *http.Request, id string) {
	var input struct {
		Done bool `json:"done"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	var t Task
	err := h.db.QueryRowContext(r.Context(),
		`UPDATE tasks SET done=$1 WHERE id=$2 RETURNING id, title, done, created_at`,
		input.Done, id,
	).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request, id string) {
	res, err := h.db.ExecContext(r.Context(), `DELETE FROM tasks WHERE id=$1`, id)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
