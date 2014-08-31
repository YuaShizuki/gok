package main
import "fmt"
import "github.com/howeyc/fsnotify"
import "io/ioutil"
import "path/filepath"
import "os/exec"
import "net"
import "strings"

var controllerListener net.Listener

func runner() {
    go startNotifier()
    err := run()
    if err != nil {
        fmt.Println(err)
    }
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

/*- http server (child process) controller -*/
func controller() (string, error) {
    controllerListener, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil {
        return "", err
    }
    port := strings.Split(l.Addr().String, ":")[1]
    go controllerStart()
    return port, nil
}

func switchOffContorler() {
    if controllerListener != nil {
        controllerListener.Close()
        controllerListener = nil
    }
}

func controllerStart() {
    conn, err := controllerListener.Accept()
    if err != nil {
        return
    }
    ioutil.ReadAll(conn)
    switchOffControler()
}
/*- End -*/


func startNotifier() {
    goorgok := regexp.Compile("(.*\\.go|.*\\.gok)$")
    watch, err := fsnotify.NewWatcher()
    if err != nil {
        fmt.Println(err)
        os.Exit()
    }
    for ;; {
        select {
            case event := <-watch.Event:
                if goorgok.Match([]byte(event.Name())) {
                    txtserve.StopServer()
                    switchOffControler()
                    err := run()
                    if err != nil {
                        txtserve.StartServer(err.Error())
                    }
                }
            case err := <-watch.Error:
                fmt.Println("runner out of controll")
                os.Exit(1)
        }
    }
}

