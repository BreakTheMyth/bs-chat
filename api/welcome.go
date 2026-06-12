package api

import (
	"net/http"
	"io/fs"
	"fmt"
	"mygo/server"
	"mygo/web"
	"mygo/log"
)

const WELCOME_METHOD string = "GET"
const WELCOME_PATH   string = "/welcome/"

func init() {
	server.Register(fmt.Sprintf("%s %s", 
		WELCOME_METHOD, WELCOME_PATH), welcome_handler)
}

func welcome_handler(w http.ResponseWriter, r *http.Request) {
	subFS, err := fs.Sub(web.BuildFS, "build")
	if err != nil {
		log.Err <- err
	}

	fileServer := http.FileServer(http.FS(subFS))

	http.StripPrefix(WELCOME_PATH, fileServer).ServeHTTP(w, r)
}
