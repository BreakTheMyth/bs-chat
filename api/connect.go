package api

import (
	"encoding/base64"
	"crypto/sha1"
	"net/http"
	"fmt"
	"mygo/server"
	"mygo/log"
)

const CONNECT_METHOD string = "GET"
const CONNECT_PATH   string = "/connect"

func init() {
	server.Register(fmt.Sprintf("%s %s",
	CONNECT_METHOD, CONNECT_PATH), connect_handler)
}

func connect_handler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Connection") != "Upgrade" ||
		r.Header.Get("Upgrade") != "websocket" {

		log.Err <- "connect_handler: Not WebSocket."
		return
	}

	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		log.Err <- "connect_handler: No Sec-WebSocker-Key."
		return
	}

	accept := computeAccept(key)

	w.Header().Set("Upgrade", "websocket")
	w.Header().Set("Connection", "Upgrade")
	w.Header().Set("Sec-WebSocket-Accept", accept)

	w.WriteHeader(101)
}

func computeAccept(key string) string {
	const magicGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

	data := key + magicGUID
	hash := sha1.Sum([]byte(data))

	return base64.StdEncoding.EncodeToString(hash[:])
}
