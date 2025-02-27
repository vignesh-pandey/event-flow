package api


import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock handler for GetUsers
func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Simulate fetching users from Redis or a database
	users := []map[string]interface{}{
		{"id": 1, "name": "Jack", "age": 30},
		{"id": 2, "name": "Raju", "age": 25},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func TestGetUsers(t *testing.T) {
	// Create an HTTP request
	req := httptest.NewRequest("GET", "/getusers", nil)

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler
	GetUsers(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	// Assert the response body
	var users []map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &users)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "Jack", users[0]["name"])
	assert.Equal(t, 30.0, users[0]["age"])
}
