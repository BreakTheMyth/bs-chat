package api

import (
    "net/http"
    "errors"
    "bytes"
    "fmt"
    "mygo/server"
    "mygo/log"
)

func init() {
    server.Register("GET /connect", connect_handler)
}

func connect_handler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()

    nickname := query.Get("nickname")
    if nickname == "" {
        w.WriteHeader(400)
        return
    }

    headshot := query.Get("headshot")
    if headshot == "" {
        w.WriteHeader(400)
        return
    }

    var ico uint32

    n, err := fmt.Sscanf(headshot, "%d", &ico)
    if n != 1 || err != nil {
        w.WriteHeader(400)
        return
    }

    user_id := user_create(nickname, ico)
    defer user_destroy(user_id)

    message := user_list[user_id].message

    in, out, done, once, err := websocket_upgrade(w, r)
    if err != nil {
        log.Err <- err
    }

    defer once.Do(func() {
        close(done)
    })

    id        := []byte(fmt.Sprintf("%d", user_id))
    header, _ := websocket_make("text", uint64(len(id)))
    out <- append(header, id...)

    for { select {
        case event := <-in:

            err = handle_event(user_id, event, out)
            if err != nil {
                close_, _ := websocket_make("close", 0)
                out <- close_
                return
            }

        case msg := <- message:

            header, _ := websocket_make("text", uint64(len(msg)))
            out <- append(header, msg...)

        case <-done: return
    } }
}

func handle_event(user_id uint32, event []byte, out chan []byte) error {
    var act        string
    var session_id uint32
    var msg        []byte
            
    _, is_exist := user_list[user_id]
    if !is_exist {
        return errors.New("invalid user_id")
    }

    parts := bytes.SplitN(event, []byte{0x1e}, 3)

    act = string(parts[0])

    n, err := fmt.Sscanf(string(parts[1]), "%d", &session_id)
    if n != 1 || err != nil {
        return errors.New("invalid session_id")
    }

    msg = parts[2]

    switch act {
        case "join":
            if len(msg) != 0 {
                session_id = session_create(string(msg))
            }

            user_join_session(user_id, session_id)

            msg = []byte(fmt.Sprintf("%d", session_id))

            header, _ := websocket_make("text", uint64(len(msg)))
            out <- append(header, msg...)

        case "exit":
        
            is_exist = user_is_in_session(user_id, session_id)
            if !is_exist {
                return errors.New("already exit")
            }

            user_exit_session(user_id, session_id)

            s := session_list[session_id]

            if len(s.user_id_list) == 0 {
                session_destroy(session_id)
            }

            header, _ := websocket_make("text", 1)
            out <- append(header, 'y')

        case "send":

            is_exist = user_is_in_session(user_id, session_id)
            if !is_exist {
                return errors.New("already exit")
            }

            message := []byte(fmt.Sprintf("%d\x1f%d\x1f", session_id, user_id))
            message = append(message, msg...)

            s := session_list[session_id]
            for k, _ := range s.user_id_list {
                if k == user_id {
                    header, _ := websocket_make("text", 1)
                    out <- append(header, 'y')
                    continue
                }

                user_list[k].message <- message
            }
    }

    return nil
}
