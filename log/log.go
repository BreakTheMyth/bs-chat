package log

import (
    "time"
    "fmt"
    "os"
)

var Info = make(chan any, 0x1000)
var Err  = make(chan any, 0x1000)

func Start() {
    go func() { for { select {

        case msg := <-Info:
            fmt.Fprintf(os.Stdout, 
                "[%v]: \033[34m%s\033[0m\n", 
                time.Now().Format("2006-01-02 15:04:05"), msg)

        case msg := <-Err:
            fmt.Fprintf(os.Stderr, 
                "[%v]: \033[31m%s\033[0m\n", 
                time.Now().Format("2006-01-02 15:04:05"), msg)

    } } }()
}
