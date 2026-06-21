package server

import (
    "net/http"
    "fmt"
    "github.com/BreakTheMyth/bs-chat/config"
    "github.com/BreakTheMyth/bs-chat/log"
)

type handler struct {
    rule     string
    function http.HandlerFunc
}

var handlerList = make([]handler, 0)

func Start() {
    mux := http.NewServeMux()

    for _, v := range handlerList {
        mux.HandleFunc(v.rule, v.function)
    }

    log.Start()

    log.Info <- fmt.Sprintf("%s:%s",
        config.SERVER_HOST, config.SERVER_PORT)

    log.Err <- http.ListenAndServeTLS(fmt.Sprintf("%s:%s",
        config.SERVER_HOST, config.SERVER_PORT),
        config.SERVER_CRT, config.SERVER_KEY, mux)
}

func Register(rule string, function http.HandlerFunc) {
    handlerList = append(handlerList, handler{rule, function})
}
