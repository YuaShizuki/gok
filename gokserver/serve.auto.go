package main
import "fmt"
import "net/http"
import "io/ioutil"
import "os"
import "net"
import "regexp"
import "strings"

var coreListener net.Listener

var routes map[string]func(*Gok) = map[string]func(*Gok) {
//<gok inject routes>
};

type mainHandler struct{};
func (_ *mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    gok := new(Go)
    gok.r = r
    gok.w = w
    gok.response = new(bytes.Buffer)
    fn, ok := routes[gok.ServerSelf()];
    if ok {
        fn(gok);
        if gok.deadMsg != nil {
            fmt.Fprintln(w, gok.deadMsg)
        } else if gok.shouldRedirect {
            fmt.Fprintln(w, "")
        }
    } else {
        if exist,_ := pathExist(gok.ServerSelf()); exist {
            http.ServeFile(w, r, gok.ServerSelf());
            return
        }
        w.WriteHeader(404);
        content404, err := ioutil.ReadFile("404.html");
        if err != nil {
            fmt.Fprintln(w, "Error Not Found!");
        } else {
            fmt.Fprintln(w, string(content404));
        }
    }
}

func main() {
    var err error
    coreListener, err = net.Listen("tcp", ":80")
    if err != nil {
        errExit(err)
    }
    go controller()
    http.Serve(coreListener, new(mainHandler))
}

func controller() {
    regsrch,_ := regexp.Compile("^gokcontroller=[0-9]+$")
    for _, command := range os.Args {
        if regsrch.Match([]byte(command)) {
            port := strings.Split(command, "=")[1]
            conn, err := net.Dial("tcp", "127.0.0.1:"+port)
            if err != nil {
                errExit(err)
            }
            defer conn.Close()
            ioutil.ReadAll(conn)
            coreListener.Close()
            return
        }
    }
}
