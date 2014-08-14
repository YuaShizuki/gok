package main
import "fmt"
import "net/http"

func index(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello World!"));
}

func main() {
    http.HandleFunc("/", index)
    err := http.ListenAndServe(":80", nil);
    if err != nil {
        fmt.Println(err);
    }
}
