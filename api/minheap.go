package api

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
