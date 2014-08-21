package main
import "fmt"
import "net/http"
import "io/ioutil"

var routes map[string]func(*Gok) = map[string]func(*Gok) {
//<gok routes>
//</gok>
};

type mainHandler struct{};
func (_ *mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    gok := new(Gok);
    gok.r = r;
    gok.w = w;
    fn, ok := routes[r.URL.Path];
    if ok {
        fn(gok);
        fmt.Fprint(w, gok.response);
    } else {
        w.WriteHeader(404);
        content404, err := ioutil.ReadFile("404.html");
        if err != nil {
            fmt.Fprintln(w, "Error Not Found!");
        } else {
            fmt.Fprintln(w, content404);
        }
    }
}

func main() {
    http.Handle("/", new(mainHandler))
    fmt.Println("server running like a bitch!");
    err := http.ListenAndServe(":8080", nil);
    if err != nil {
        fmt.Println(err);
    }
}
