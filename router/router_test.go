package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/websocket"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	// Set the working directory to the project root
	err := os.Chdir("..") // Adjust the path as needed
	if err != nil {
		t.Fatalf("Failed to set working directory: %v", err)
	}

	// Create a new router
	r := chi.NewRouter()
	SetupRouter(r)

	t.Run("TestCreateRoom", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/rooms", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status OK for /api/rooms")
	})

	t.Run("TestCheckRoom", func(t *testing.T) {
		// Create a room first to ensure it exists
		roomID := createRoomAndGetId(r)
		req := httptest.NewRequest(http.MethodGet, "/api/rooms/"+roomID, nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status OK for /api/rooms/{id}")
	})

	t.Run("TestJoinRoom", func(t *testing.T) {
		// Create a room first to ensure it exists
		roomID := createRoomAndGetId(r)
		t.Logf("Room ID: %s", roomID)

		// Create a test server using the router
		server := httptest.NewServer(r)
		defer server.Close()

		// Convert the test server URL to a WebSocket URL
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/rooms/" + roomID + "/join"
		t.Logf("WebSocket URL: %s", wsURL)

		// Connect to the WebSocket server
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("Failed to connect to WebSocket: %v", err)
		}
		defer ws.Close()

		// Verify the WebSocket connection
		assert.NotNil(t, ws, "WebSocket connection should not be nil")
	})

	t.Run("TestServeIndex", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status OK for /")
		assert.Contains(t, rec.Body.String(), "<!DOCTYPE html>", "Expected HTML content for /")
	})

	t.Run("TestServeStaticFiles", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/static/script.js", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status OK for /static/*")
		assert.True(t, strings.Contains(rec.Header().Get("Content-Type"), "text/javascript"), "Expected CSS content type for /static/*")
	})
}

func createRoomAndGetId(r *chi.Mux) string {
	createReq := httptest.NewRequest(http.MethodPost, "/api/rooms", nil)
	createRec := httptest.NewRecorder()
	r.ServeHTTP(createRec, createReq)
	var response map[string]string
	json.Unmarshal(createRec.Body.Bytes(), &response)
	roomID, _ := response["room_id"]
	return roomID
}
