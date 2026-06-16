package api

import (
    "container/heap"
    "sync"
    "mygo/config"
)

type user struct {
    user_id         uint32
    user_name       string
    user_ico        string
    message         chan string
    session_id_list map[uint32]struct{}
}

var user_count = uint32(0)
var user_list  = make(map[uint32]user)
var idle_user  = &uint32_min_heap{}

var user_mutex sync.RWMutex

func init() {
    heap.Init(idle_user)
}

func user_create(user_name string, user_ico string) uint32 {
    var user_id uint32

    user_mutex.Lock()
    defer user_mutex.Unlock()

    user_id = user_count
    if idle_user.Len() != 0 {
        user_id = heap.Pop(idle_user).(uint32)
    }

    user_count++

    user_list[user_id] = user{
        user_id:         user_id, 
        user_name:       user_name, 
        user_ico:        "1", 
        message:         make(chan string, config.BUFFER_SIZE),
        session_id_list: map[uint32]struct{}{},
    }

    return user_id
}

func user_destroy(user_id uint32) {
    user_mutex.Lock()
    defer user_mutex.Unlock()
    session_mutex.Lock()
    defer session_mutex.Unlock()

    user_count--

    u := user_list[user_id]

    for session_id := range u.session_id_list {
        s := session_list[session_id]

        delete(s.user_id_list, user_id)
    }

    u.message         = nil
    u.session_id_list = nil

    delete(user_list, user_id)

    heap.Push(idle_user, user_id)
}

func user_is_in_session(user_id uint32, session_id uint32) bool {
    user_mutex.RLock()
    defer user_mutex.RUnlock()

    var is_exist bool
    var u        user

    u, is_exist = user_list[user_id]
    if !is_exist {
        return false
    }

    _, is_exist = u.session_id_list[session_id]

    return is_exist
}

func user_join_session(user_id uint32, session_id uint32) {
    user_mutex.Lock()
    defer user_mutex.Unlock()
    session_mutex.Lock()
    defer session_mutex.Unlock()

    u := user_list[user_id]
    s := session_list[session_id]

    u.session_id_list[session_id] = struct{}{}
    s.user_id_list[user_id]       = struct{}{}
}

func user_exit_session(user_id uint32, session_id uint32) {
    user_mutex.Lock()
    defer user_mutex.Unlock()
    session_mutex.Lock()
    defer session_mutex.Unlock()

    u := user_list[user_id]
    s := session_list[session_id]

    delete(u.session_id_list, session_id)
    delete(s.user_id_list, user_id)
}
