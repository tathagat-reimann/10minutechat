package websocketX

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckOrigin_AllowedHost(t *testing.T) {
	// Set the ALLOWED_HOST environment variable
	os.Setenv("ALLOWED_HOST", "example.com")
	defer os.Unsetenv("ALLOWED_HOST")

	// Create a mock HTTP request with the allowed host
	req := &http.Request{Host: "example.com"}

	// Check if the origin is allowed
	isAllowed := Upgrader.CheckOrigin(req)

	assert.True(t, isAllowed, "Expected CheckOrigin to allow the request from the allowed host")
}

func TestCheckOrigin_DisallowedHost(t *testing.T) {
	// Set the ALLOWED_HOST environment variable
	os.Setenv("ALLOWED_HOST", "example.com")
	defer os.Unsetenv("ALLOWED_HOST")

	// Create a mock HTTP request with a disallowed host
	req := &http.Request{Host: "notallowed.com"}

	// Check if the origin is disallowed
	isAllowed := Upgrader.CheckOrigin(req)

	assert.False(t, isAllowed, "Expected CheckOrigin to disallow the request from a disallowed host")
}

func TestCheckOrigin_DefaultHost(t *testing.T) {
	// Unset the ALLOWED_HOST environment variable to use the default
	os.Unsetenv("ALLOWED_HOST")

	// Create a mock HTTP request with the default host
	req := &http.Request{Host: "localhost:8080"}

	// Check if the origin is allowed
	isAllowed := Upgrader.CheckOrigin(req)

	assert.True(t, isAllowed, "Expected CheckOrigin to allow the request from the default host")
}

func TestMaxRoomCapacity_ValidEnv(t *testing.T) {
	// Set a valid MAX_ROOM_CAPACITY environment variable
	os.Setenv("MAX_ROOM_CAPACITY", "5")
	defer os.Unsetenv("MAX_ROOM_CAPACITY")

	// Verify that MaxRoomCapacity is set correctly
	assert.Equal(t, 5, MaxRoomCapacity(), "Expected MaxRoomCapacity to be 5")
}

func TestMaxRoomCapacity_InvalidEnv(t *testing.T) {
	// Set an invalid MAX_ROOM_CAPACITY environment variable
	os.Setenv("MAX_ROOM_CAPACITY", "invalid")
	defer os.Unsetenv("MAX_ROOM_CAPACITY")

	// Verify that MaxRoomCapacity falls back to the default value
	assert.Equal(t, 2, MaxRoomCapacity(), "Expected MaxRoomCapacity to default to 2")
}

func TestMaxRoomCapacity_Default(t *testing.T) {
	// Unset the MAX_ROOM_CAPACITY environment variable
	os.Unsetenv("MAX_ROOM_CAPACITY")

	// Verify that MaxRoomCapacity defaults to 2
	assert.Equal(t, 2, MaxRoomCapacity(), "Expected MaxRoomCapacity to default to 2")
}
