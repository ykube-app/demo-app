package api_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ykube-app/demo-app/internal/api"
	"github.com/ykube-app/demo-app/internal/db"
)

func testDB(t *testing.T) *sql.DB {
	t.Helper()
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL environment variable is not set")
	}
	pool, err := db.Open()
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.Exec("DELETE FROM tasks")
		pool.Close()
	})
	return pool
}

func TestListTasks_Empty(t *testing.T) {
	pool := testDB(t)
	h := api.NewHandler(pool)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var tasks []map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &tasks))
	assert.Empty(t, tasks)
}

func TestCreateAndListTask(t *testing.T) {
	pool := testDB(t)
	h := api.NewHandler(pool)

	body := `{"title": "Buy Milk"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var created map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &created))
	assert.Equal(t, "Buy Milk", created["title"])
	assert.Equal(t, false, created["done"])
	id := created["id"].(string)

	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	h.ServeHTTP(rec2, req2)

	var tasks []map[string]interface{}
	require.NoError(t, json.Unmarshal(rec2.Body.Bytes(), &tasks))
	assert.Len(t, tasks, 1)
	assert.Equal(t, id, tasks[0]["id"])
}

func TestToggleTask(t *testing.T) {
	pool := testDB(t)
	h := api.NewHandler(pool)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(`{"title": "Buy Milk"}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var created map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &created))
	id := created["id"].(string)

	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPatch, "/api/tasks/"+id, bytes.NewBufferString(`{"done": true}`))
	req2.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rec2, req2)
	assert.Equal(t, http.StatusOK, rec2.Code)

	var updated map[string]interface{}
	require.NoError(t, json.Unmarshal(rec2.Body.Bytes(), &updated))
	assert.Equal(t, true, updated["done"])
}

func TestDeleteTask(t *testing.T) {
	pool := testDB(t)
	h := api.NewHandler(pool)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(`{"title": "Buy Milk"}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var created map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &created))
	id := created["id"].(string)

	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodDelete, "/api/tasks/"+id, nil)
	h.ServeHTTP(rec2, req2)
	assert.Equal(t, http.StatusNoContent, rec2.Code)
}
