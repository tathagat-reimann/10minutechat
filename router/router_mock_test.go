package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// Mock handlers for the room package
func mockCreateRoom(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"room_id": "mockRoom123"}`))
}

func mockCheckRoom(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"room_id": "mockRoom123"}`))
}

func mockJoinRoom(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusSwitchingProtocols)
}

// Test class for router with mocked room handlers
func TestRouterWithMockRoom(t *testing.T) {
	// Create a new router
	r := chi.NewRouter()

	// Replace room handlers with mock handlers
	r.Route("/api", func(r chi.Router) {
		r.Post("/rooms", mockCreateRoom)
		r.Get("/rooms/{id}", mockCheckRoom)
		r.Get("/rooms/{id}/join", mockJoinRoom)
	})

	t.Run("TestMockCreateRoom", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/rooms", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status OK for /api/rooms")
		assert.Contains(t, rec.Body.String(), `"room_id": "mockRoom123"`, "Expected mock room_id in response")
	})

	t.Run("TestMockCheckRoom", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/rooms/mockRoom123", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status OK for /api/rooms/{id}")
		assert.Contains(t, rec.Body.String(), `"room_id": "mockRoom123"`, "Expected mock room_id in response")
	})

	t.Run("TestMockJoinRoom", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/rooms/mockRoom123/join", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusSwitchingProtocols, rec.Code, "Expected status Switching Protocols for /api/rooms/{id}/join")
	})
}
