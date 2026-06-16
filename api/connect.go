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
    in, out, done, once, err := websocket_upgrade(w, r)
    if err != nil {
        log.Err <- err
    }

    _ = in
    _ = out
    _ = done
    _ = once

    for {
        msg := <-in

        log.Info <- msg

        if string(msg) == "Hello" {
            hi        := []byte("Hi")
            header, _ := websocket_make("text", uint64(len(hi)))
            out <- append(header, hi...)
        }

        if string(msg) == "exit" {
            close_, _ := websocket_make("close", 0)
            out <- close_
            once.Do(func() {
                close(done)
            })
            return
        }
    }
}
