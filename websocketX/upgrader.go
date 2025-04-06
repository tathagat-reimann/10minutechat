package websocketX

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tathagat/10minutechat/conf"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return r.Host == conf.AllowedHost
	},
}
