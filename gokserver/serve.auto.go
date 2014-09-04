package main
import "fmt"
import "net/http"
import "io/ioutil"
import "os"
import "net"
import "regexp"
import "strings"
import "bytes"

var gok_coreListener net.Listener

var gok_routes map[string]func(*Gok) = map[string]func(*Gok) {
//<gok inject routes>
};

type gok_mainHandler struct{};
func (_ *gok_mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    gok := buildGok(w, r)
    fn, ok := gok_routes[gok.ServerSelf()];
    if ok {
        fn(gok);
        if gok.deadMsg != "" {
            fmt.Fprintln(w, gok.deadMsg)
        } else if gok.shouldRedirect {
            fmt.Fprintln(w, "")
        } else {
            fmt.Fprintln(w, gok.response.String())
        }
    } else if !handleIfQuickAjx(gok) {
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
    gok_coreListener, err = net.Listen("tcp", ":80")
    if err != nil {
        errExit(err)
    }
    go controller()
    http.Serve(gok_coreListener, new(gok_mainHandler))
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
            gok_coreListener.Close()
            return
        }
    }
}

func buildGok(w http.ResponseWriter, r *http.Request) *Gok {
    result := new(Gok)
    result.w = w
    result.r = r
    result.deadMsg = ""
    result.response = new(bytes.Buffer)
    result.shouldRedirect = false
    return result
}
