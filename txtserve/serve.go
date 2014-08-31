package main
import "fmt"
import "net"
import "net/http"
import "errors"
import "time"

var running bool = false
var onlyResponse string
var mainListener net.Listener

type mainHandler struct {}
func (_ *mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, time.Now().Unix())
}

func serv() {
    if err := http.Serve(mainListener, new(mainHandler)); err != nil {
        StopServer()
    }
}

func StartServer(content string) error {
    l, err := net.Listen("tcp", ":80")
    if err != nil {
        return err
    }
    mainListener = l
    onlyResponse = content
    running = true
    go serv()
    return nil
}

func StopServer() error {
    running = false
    if mainListener == nil {
        return errors.New("server not listening")
    }
    err := mainListener.Close()
    return err
}

func IsRunning() bool {
    return running
}
