package log

import (
    "time"
    "fmt"
    "os"
    "mygo/config"
)

var Info = make(chan any, config.BUFFER_SIZE)
var Err  = make(chan any, config.BUFFER_SIZE)

func Start() {
    go func() { for { select {

        case msg := <-Info:
            fmt.Fprintf(os.Stdout, 
                "\n> %v\n\033[34m%s\033[0m\n", 
                time.Now(), msg)

        case msg := <-Err:
            fmt.Fprintf(os.Stderr, 
                "\n> %v\n\033[31m%s\033[0m\n", 
                time.Now(), msg)

    } } }()
}
