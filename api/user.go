package api

import (
    "container/heap"
    "sync"
)

type user struct {
    user_id         uint32
    user_name       string
    user_ico        string
    session_id_list map[uint32]struct{}
}

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

var user_count = uint32(0)
var user_list  = make(map[uint32]user)
var idle_user  = &uint32_min_heap{}
var mutex sync.Mutex

func init() {
    heap.Init(idle_user)
}

func user_create(user_name string) uint32 {
    var user_id uint32

    mutex.Lock()

    user_id = user_count
    if idle_user.Len() != 0 {
        user_id = heap.Pop(idle_user).(uint32)
    }

    user_count++

    mutex.Unlock()

    user_list[user_id] = user{
        user_id:         user_id, 
        user_name:       user_name, 
        user_ico:        "1", 
        session_id_list: map[uint32]struct{}{},
    }

    return user_id
}

func user_destroy(user_id uint32) {
    mutex.Lock()

    user_count--

    delete(user_list, user_id)

    heap.Push(idle_user, user_id)

    mutex.Unlock()
}
