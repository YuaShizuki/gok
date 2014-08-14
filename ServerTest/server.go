package main
import "fmt"
import "net/http"

type Gok struct {
}

type mainHandler struct{};
func (_ *mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "<html>", "Javc", 32, "</html>");
}

func main() {
    http.Handle("/", new(mainHandler))
    err := http.ListenAndServe(":8080", nil);
    if err != nil {
        fmt.Println(err);
    }
}
