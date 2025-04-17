package room

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/tathagat/10minutechat/conf"
)

// func TestJoinRoom(t *testing.T) {
// 	t.Run("JoinRoom_Success", func(t *testing.T) {
// 		// Create a room
// 		roomID := CreateRoom()

// 		// Mock a WebSocket connection
// 		conn := &websocket.Conn{}

// 		// Call JoinRoom
// 		JoinRoom(roomID, conn)

// 		// Verify the client was added to the room
// 		roomsMu.Lock()
// 		room, exists := rooms[roomID]
// 		roomsMu.Unlock()

// 		if !exists {
// 			t.Fatalf("Room with ID %s does not exist", roomID)
// 		}

// 		room.Mutex.Lock()
// 		clientName, clientExists := room.Clients[conn]
// 		room.Mutex.Unlock()

// 		if !clientExists {
// 			t.Errorf("Client was not added to the room")
// 		}

// 		if clientName == "" {
// 			t.Errorf("Client name was not assigned")
// 		}
// 	})

// 	t.Run("JoinRoom_RoomDoesNotExist", func(t *testing.T) {
// 		// Mock a WebSocket connection
// 		conn := &websocket.Conn{}

// 		// Call JoinRoom with a non-existent room ID
// 		JoinRoom("nonexistent-room-id", conn)

// 		// Verify no room was created
// 		roomsMu.Lock()
// 		_, exists := rooms["nonexistent-room-id"]
// 		roomsMu.Unlock()

// 		if exists {
// 			t.Errorf("Room with ID 'nonexistent-room-id' should not exist")
// 		}
// 	})

// 	t.Run("JoinRoom_RoomFull", func(t *testing.T) {
// 		// Set MaxRoomCapacity to 1 for testing
// 		conf.MaxRoomCapacity = 1

// 		// Create a room
// 		roomID := CreateRoom()

// 		// Mock two WebSocket connections
// 		conn1 := &websocket.Conn{}
// 		conn2 := &websocket.Conn{}

// 		// Add the first client to the room
// 		JoinRoom(roomID, conn1)

// 		// Attempt to add the second client to the room
// 		JoinRoom(roomID, conn2)

// 		// Verify the second client was not added
// 		roomsMu.Lock()
// 		room, exists := rooms[roomID]
// 		roomsMu.Unlock()

// 		if !exists {
// 			t.Fatalf("Room with ID %s does not exist", roomID)
// 		}

// 		room.Mutex.Lock()
// 		_, client2Exists := room.Clients[conn2]
// 		room.Mutex.Unlock()

// 		if client2Exists {
// 			t.Errorf("Second client should not have been added to the room as it is full")
// 		}
// 	})
// }

func TestCreateRoom(t *testing.T) {
	t.Run("CreateRoom_Success", func(t *testing.T) {
		// Create a room
		roomID := CreateRoom()

		// Verify the room was created
		roomsMu.Lock()
		room, exists := rooms[roomID]
		roomsMu.Unlock()

		if !exists {
			t.Fatalf("Room with ID %s does not exist", roomID)
		}

		if room.ID != roomID {
			t.Errorf("Expected room ID to be %s, got %s", roomID, room.ID)
		}

		if room.Broadcast == nil {
			t.Errorf("Broadcast channel was not initialized")
		}

		if len(room.Clients) != 0 {
			t.Errorf("Expected no clients in the room, but found %d", len(room.Clients))
		}
	})
}

func TestCheckRoom(t *testing.T) {
	t.Run("CheckRoom_Success_Found", func(t *testing.T) {
		// Create a room
		roomID := CreateRoom()

		// Test
		exists, _, _ := CheckRoom(roomID)
		if !exists {
			t.Errorf("Expected room to exist, but it does not")
		}
	})
	t.Run("CheckRoom_Success_NotFound", func(t *testing.T) {
		// Create a room
		roomID := "fake-room-id"

		// Test
		exists, _, _ := CheckRoom(roomID)
		if exists {
			t.Errorf("Expected room to NOT exist, but it does ")
		}
	})
	t.Run("CheckRoom_Success_Full", func(t *testing.T) {
		// Set MaxRoomCapacity to 1 for testing
		conf.MaxRoomCapacity = 1
		// Create a room
		roomID := CreateRoom()
		// Mock a WebSocket connection
		conn := &websocket.Conn{}
		// Call JoinRoom
		JoinRoom(roomID, conn)

		// Test
		_, full, _ := CheckRoom(roomID)
		if !full {
			t.Errorf("Expected room to be full, but it is not")
		}
	})
}
