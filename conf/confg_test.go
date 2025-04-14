package conf

import (
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestInit(t *testing.T) {
	// Save original environment variables to restore later
	originalAllowedHost := os.Getenv("ALLOWED_HOST")
	originalMaxRoomCapacity := os.Getenv("MAX_ROOM_CAPACITY")

	// Restore environment variables after tests
	defer func() {
		os.Setenv("ALLOWED_HOST", originalAllowedHost)
		os.Setenv("MAX_ROOM_CAPACITY", originalMaxRoomCapacity)
	}()

	t.Run("DefaultAllowedHost", func(t *testing.T) {
		os.Unsetenv("ALLOWED_HOST")
		if AllowedHost != "localhost:8080" {
			t.Errorf("Expected AllowedHost to be 'localhost:8080', got '%s'", AllowedHost)
		}
	})

	t.Run("CustomAllowedHost", func(t *testing.T) {
		os.Setenv("ALLOWED_HOST", "example.com:9000")
		LoadConfig()
		if AllowedHost != "example.com:9000" {
			t.Errorf("Expected AllowedHost to be 'example.com:9000', got '%s'", AllowedHost)
		}
	})

	t.Run("DefaultMaxRoomCapacity", func(t *testing.T) {
		os.Unsetenv("MAX_ROOM_CAPACITY")
		if MaxRoomCapacity != 2 {
			t.Errorf("Expected MaxRoomCapacity to be 2, got %d", MaxRoomCapacity)
		}
	})

	t.Run("CustomMaxRoomCapacity", func(t *testing.T) {
		os.Setenv("MAX_ROOM_CAPACITY", "5")
		LoadConfig()
		if MaxRoomCapacity != 5 {
			t.Errorf("Expected MaxRoomCapacity to be 5, got %d", MaxRoomCapacity)
		}
	})

	t.Run("InvalidMaxRoomCapacity", func(t *testing.T) {
		os.Setenv("MAX_ROOM_CAPACITY", "invalid")
		LoadConfig()
		if MaxRoomCapacity != 2 {
			t.Errorf("Expected MaxRoomCapacity to fall back to 2, got %d", MaxRoomCapacity)
		}
	})

	t.Run("RandomNamesLengthCheck", func(t *testing.T) {
		os.Setenv("MAX_ROOM_CAPACITY", strconv.Itoa(len(RandomNames)+1))

		if os.Getenv("BE_CRASHER") == "1" {
			LoadConfig()
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=RandomNamesLengthCheck")
		cmd.Env = append(os.Environ(), "BE_CRASHER=1")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		t.Fatalf("process ran with err %v, want exit status 1", err)
	})
}
