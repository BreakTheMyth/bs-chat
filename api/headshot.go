package api

import (
    "net/http"
    "io/fs"
    "github.com/BreakTheMyth/bs-chat/server"
    "github.com/BreakTheMyth/bs-chat/web"
    "github.com/BreakTheMyth/bs-chat/log"
)

func init() {
    server.Register("GET /headshot/", headshot_handler)
}

func headshot_handler(w http.ResponseWriter, r *http.Request) {
    subFS, err := fs.Sub(web.HeadshotFS, "headshot")
    if err != nil {
        log.Err <- err
    }

    fileServer := http.FileServer(http.FS(subFS))

    http.StripPrefix("/headshot/", fileServer).ServeHTTP(w, r)
}
