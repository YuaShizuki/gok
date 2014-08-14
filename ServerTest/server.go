package main
import "fmt"
import "net/http"

type Gok struct {
    w http.ResponseWriter;
    r *http.Request;
    Url string;
};

func (gok *Gok) Echo(a ...interface{}) {}

func (gok *Gok) ServerSelf() string {}
func (gok *Gok) ServerHttpUserAgent() string {}
func (gok *Gok) ServerHttps() bool {}
func (gok *Gok) ServerRemoteAddr() string {}
func (gok *Gok) ServerRemotePort() int {}
func (gok *Gok) ServerPort() int {}
func (gok *Gok)


type handel func(*Gok);

type mainHandler struct{};
func (_ *mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "<html>", r.URL.RawQuery, "</html>");
}

func main() {
    http.Handle("/", new(mainHandler))
    err := http.ListenAndServe(":8080", nil);
    if err != nil {
        fmt.Println(err);
    }
}
