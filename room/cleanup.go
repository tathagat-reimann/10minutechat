package room

import (
    "log"
    "time"
)

func CleanupRooms() {
    for {
        time.Sleep(1 * time.Minute)
        roomsMu.Lock()
        for id, room := range rooms {
            if len(room.Clients) == 0 {
                delete(rooms, id)
                log.Printf("Room deleted: %s", id)
            }
        }
        roomsMu.Unlock()
    }
}