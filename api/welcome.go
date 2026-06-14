package api

import (
    "net/http"
    "io/fs"
    "mygo/server"
    "mygo/web"
    "mygo/log"
)

func init() {
    server.Register("GET /welcome/", welcome_handler)
}

func welcome_handler(w http.ResponseWriter, r *http.Request) {
    subFS, err := fs.Sub(web.BuildFS, "build")
    if err != nil {
        log.Err <- err
    }

    fileServer := http.FileServer(http.FS(subFS))

    http.StripPrefix("/welcome/", fileServer).ServeHTTP(w, r)
}
