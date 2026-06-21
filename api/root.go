package api

import (
	"io/fs"
	"mygo/log"
	"mygo/server"
	"mygo/web"
	"net/http"
)

func init() {
    server.Register("GET /", root_handler)
}

func root_handler(w http.ResponseWriter, r *http.Request) {
    subFS, err := fs.Sub(web.BuildFS, "build")
    if err != nil {
        log.Err <- err
    }

    fileServer := http.FileServer(http.FS(subFS))

    fileServer.ServeHTTP(w, r)
}
