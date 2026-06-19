package api

import (
    "container/heap"
    "net/http"
    "sync"
    "fmt"
    "mygo/server"
    "mygo/config"
)

type session struct {
    session_id   uint32
    session_name string
    user_id_list map[uint32]struct{}
}

var session_count = uint32(0)
var session_list  = make(map[uint32]session)
var idle_session  = &uint32_min_heap{}

var session_mutex sync.Mutex

func init() {
    heap.Init(idle_session)

    server.Register("GET /sessions/{id}", sessions_handler)
    server.Register("GET /sessions/list", sessions_list_handler)
}

func sessions_handler(w http.ResponseWriter, r *http.Request) {
    var session_id uint32

    n, err := fmt.Sscanf(r.PathValue("id"), "%d", &session_id)
    if n != 1 || err != nil {
        w.WriteHeader(400)
        return
    }

    s, is_exist := session_list[session_id]
    if !is_exist {
        w.WriteHeader(404)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(
        fmt.Sprintf("{\"theme\":\"%s\",\"online\":%d}",
        s.session_name, len(s.user_id_list))))
}

func sessions_list_handler(w http.ResponseWriter, r *http.Request) {
    count := uint32(len(session_list))

    if count > config.LIST_MAX {
        count = config.LIST_MAX
    }

    w.Header().Set("Content-Type", "application/json")

    resp := "["

    for _, s := range session_list {
        resp = fmt.Sprintf("%s{\"id\":%d,\"theme\":\"%s\",\"online\":%d}", 
            resp, s.session_id, s.session_name, len(s.user_id_list))

        count--

        if count == 0 {
            break
        }

        resp += ","
    }

    resp += "]"

    w.Write([]byte(resp))
}

func session_create(session_name string) uint32 {
    session_mutex.Lock()
    defer session_mutex.Unlock()

    var session_id uint32

    session_id = session_count
    if idle_session.Len() != 0 {
        session_id = heap.Pop(idle_session).(uint32)
    }

    session_count++

    session_list[session_id] = session{
        session_id:   session_id, 
        session_name: session_name, 
        user_id_list: map[uint32]struct{}{},
    }

    return session_id
}

func session_destroy(session_id uint32) {
    user_mutex.Lock()
    defer user_mutex.Unlock()
    session_mutex.Lock()
    defer session_mutex.Unlock()

    session_count--

    s := session_list[session_id]

    for user_id := range s.user_id_list {
        u := user_list[user_id]

        delete(u.session_id_list, session_id)
    }

    s.user_id_list = nil

    delete(session_list, session_id)

    heap.Push(idle_session, session_id)
}
