package main
import "fmt"
import "github.com/howeyc/fsnotify"
import "io/ioutil"
import "path/filepath"
import "os/exec"
import "os"
import "net"
import "strings"
import "regexp"
import "./txtserve"

var controllerListener net.Listener
var controllerConn net.Conn

func runner() {
    err := run()
    if err != nil {
        errExit(err, "")
    }
    done := make(chan bool)
    <-done
}

func run() error {
    err := build(false)
    if err != nil {
        return err
    }
    pwd,_ := os.Getwd()
    exe := filepath.Base(pwd)
    port, err := controller()
    if err != nil {
        return err
    }
    cmd := exec.Command("sudo", "./"+exe, "gokcontroller="+port)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    go cmd.Run()
    return nil
}

/*- http server (child process) controller -*/
func controller() (string, error) {
    var err error
    controllerListener, err = net.Listen("tcp", "127.0.0.1:0")
    if err != nil {
        return "", err
    }
    port := strings.Split(controllerListener.Addr().String(), ":")[1]
    go controllerStart()
    return port, nil
}

func switchOffController() {
    if controllerListener != nil {
        controllerListener.Close()
        controllerListener = nil
    }
    if controllerConn != nil {
        controllerConn.Close()
        controllerConn = nil
    }
}

func controllerStart() {
    var err error
    controllerConn, err = controllerListener.Accept()
    if err != nil {
        return
    }
    go startNotifier()
    ioutil.ReadAll(controllerConn)
    switchOffController()
}
/*- End -*/


func startNotifier() {
    goorgok,_ := regexp.Compile("(.*\\.go|.*\\.gok)$")
    watch, err := fsnotify.NewWatcher()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    watch.Watch(".")
    for ;; {
        select {
            case event := <-watch.Event:
                if goorgok.Match([]byte(event.Name)) {
                    txtserve.StopServer()
                    switchOffController()
                    waitTillPortShutDown(":80")
                    err := run()
                    if err != nil {
                        fmt.Println("cannot update server, error in go|gok code")
                        txtserve.StartServer(err.Error())
                        return
                    } else {
                        fmt.Println("restarting updated server")
                        watch.Close()
                        return
                    }
                }
        }
    }
}

