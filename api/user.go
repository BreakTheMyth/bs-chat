package api

import (
    "container/heap"
    "encoding/json"
    "net/http"
    "mygo/server"
)

type User struct {
    UserId        string          `json:"user_id"`
    UserName      string          `json:"user_name"`
    UserIco       string          `json:"user_ico"`
    SessionIdList map[uint32]bool `json:"session_id_list"`
}

var user_list = map[string]User{}

type uint32_min_heap []uint32

func (h uint32_min_heap) Len() int { return len(h) }

func (h uint32_min_heap) Less(a int, b int) bool { return h[a] < h[b] }

func (h uint32_min_heap) Swap(a int, b int) { h[a], h[b] = h[b], h[a] }

func (h *uint32_min_heap) Push(i any) { *h = append(*h, i.(uint32)) }

func (h *uint32_min_heap) Pop() any {
    new_len := h.Len() - 1
    last    := (*h)[new_len]
    *h       = (*h)[:new_len]
    return last
}

func init() {

    server.Register("GET /users/{id}", get_user_handler)

    server.Register("POST /users", post_user_handler)

    server.Register("DELETE /users/{id}", delete_user_handler)

    server.Register("PATCH /users/{id}", patch_user_handler)

    idle_user := &uint32_min_heap{}

    heap.Init(idle_user)

}

func get_user_handler(w http.ResponseWriter, r *http.Request) {
    user_id := r.PathValue("id")

    user, is_exist := user_list[user_id]

    if !is_exist {
        w.WriteHeader(404)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func post_user_handler(w http.ResponseWriter, r *http.Request) {
}

func delete_user_handler(w http.ResponseWriter, r *http.Request) {
}

func patch_user_handler(w http.ResponseWriter, r *http.Request) {
}
