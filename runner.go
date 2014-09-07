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
import "time"

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
    go startNotifier(".", make(chan bool))
    ioutil.ReadAll(controllerConn)
    switchOffController()
}
/*- End -*/


func startNotifier(dir string, end chan bool) {
    lastUpdate := time.Now()
    goorgok,_ := regexp.Compile("^([^.]*\\.go|[^.]*\\.gok)$")
    watch, err := fsnotify.NewWatcher()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    dirNotifiers := make([](chan bool),0,10)
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        errExit(err, "")
    }
    for _, f := range files {
        if f.IsDir() {
            lenD := len(dirNotifiers)
            dirNotifiers = append(dirNotifiers, make(chan bool))
            startNotifier(dir+"/"+f.Name(), dirNotifiers[lenD])
        }
    }
    watch.Watch(dir)
    for ;; {
        select {
            case event := <-watch.Event:
                if goorgok.Match([]byte(event.Name)) {
                    if time.Since(lastUpdate) < (1000 * time.Millisecond) {
                        continue
                    }
                    lastUpdate = time.Now()
                    fmt.Printf("%v: src code update => ", lastUpdate)
                    txtserve.StopServer()
                    switchOffController()
                    time.Sleep(200 * time.Millisecond)
                    err = run()
                    if err != nil {
                        fmt.Printf("error building server binary\n")
                        txtserve.StartServer(err.Error())
                    } else {
                        fmt.Printf("server binary build successful\n")
                        watch.Close()
                        for _, c := range dirNotifiers {
                            c <- true
                            close(c)
                        }
                        return
                    }
                }
            case shouldEnd := <-end:
                if shouldEnd {
                    watch.Close()
                    for _, c := range dirNotifiers {
                            c <- true
                            close(c)
                    }
                    return
                }
        }
    }
}

