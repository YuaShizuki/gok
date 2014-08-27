package main
import "fmt"
import "github.com/howeyc/fsnotify"
import "io/ioutil"
import "path/filepath"
import "os/exec"
import "net"
import "strings"


var closeHttpServer func() error

func runner() {
    err := build(false)
    if err != nil {
        fmt.Println(err)
        return
    }
    pwd,_ := os.Getwd()
    exe := filepath.Base(pwd)
    port, err := controler()

    cmd := exec.Command("sudo", exe, "gokcontroller="+port)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    go notifier()
    cmd.Run()
    end := make(chan bool)
    <-end
}

func run() error {
    err := build(false)
    if err != nil {
        return err
    }
    pwd,_ := os.Getwd()
    exe := filepath.Base(pwd)
    port, err := controler()
    if err != nil {
        return err
    }
    cmd := exec.Command("sudo", exe, "gokcontroller="+port)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Run()
}

func controler() (string, error) {
    l, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil {
        return err
    }
    port := strings.Split(l.Addr().String, ":")[1]
    go holdingServer(l)
    return port, nil
}

func holdingServer(l net.Listener) {
    conn, err := l.Accept()
    closeHttpServer = conn.Close
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    ioutil.ReadAll(conn)
    defer l.Close()
}

func notifier() {
    goorgok := regexp.Compile("(.*\\.go|.*\\.gok)$")
    watch, err := fsnotify.NewWatcher()
    if err != nil {
        fmt.Println(err)
        os.Exit()
    }
    for ;; {
        select {
            case event := <-watch.Event
                if goorgok.Match([]byte(event.Name())) {
                    closeHttpServer()
                    err := run()
                    if err != nil {
                        startErrorServer(err)
                    } else {
                        shutdownErrorServer()
                    }
                }
        }
    }
}
