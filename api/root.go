package api

import (
	"net/http"
	"io/fs"
	"github.com/BreakTheMyth/bs-chat/log"
	"github.com/BreakTheMyth/bs-chat/server"
	"github.com/BreakTheMyth/bs-chat/web"
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
