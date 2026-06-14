package api

import (
    "net/http"
    "mygo/server"
    "mygo/log"
)

func init() {
    server.Register("GET /connect", connect_handler)
}

func connect_handler(w http.ResponseWriter, r *http.Request) {
    err := websocket_upgrade(w, r)
    if err != nil {
        log.Err <- err
    }
}
