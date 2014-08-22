package main
import "fmt"
import "net/http"
import "io/ioutil"

var routes map[string]func(*Gok) = map[string]func(*Gok) {
//<gok inject routes>
};

type mainHandler struct{};
func (_ *mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    gok := new(Gok);
    gok.r = r;
    gok.w = w;
    fn, ok := routes[gok.ServerSelf()];
    if ok {
        fn(gok);
        fmt.Fprint(w, gok.response);
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
    http.Handle("/", new(mainHandler))
    fmt.Println("server running like a bitch!");
    err := http.ListenAndServe(":80", nil);
    if err != nil {
        fmt.Println(err);
    }
}
