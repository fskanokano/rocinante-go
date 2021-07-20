package rocinante

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type WebsocketHandler func(*websocket.Conn)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}
