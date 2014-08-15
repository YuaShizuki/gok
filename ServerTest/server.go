package main
import "fmt"
import "net/http"
improt "flags"

/*type handel func(*Gok);
var routes map[string]handel;*/



type mainHandler struct{};
func (_ *mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    html := ```
<html>
    <body>
        <h1>Die Inside Me,...,</h1>
    </body>
</html>```;
    fmt.Fprintln(w, html);
}

func main() {
    http.Handle("/", new(mainHandler))
    err := http.ListenAndServe(":8080", nil);
    if err != nil {
        fmt.Println(err);
    }
}
