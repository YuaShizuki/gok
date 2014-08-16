package main
import "fmt"
import "net/http"
import "io"

/*type handel func(*Gok);
var routes map[string]handel;*/

func index(w http.ResponseWriter, r *http.Request) {
    var html string =
`
<html>
    <body>
        <h1>Die Inside Me!</h1>
        <form enctype='multipart/form-data' name="myGorm" action="/postit" method="POST">
            First name: <input type="text" name="firstname"><br>
            Last name: <input type="text" name="lastname">
            File: <input type="file" name="content"> 
            <input type="submit" value="Send">
        </form>
    </body>
</html>
`;
    fmt.Fprintln(w, html);
}

func postit(w http.ResponseWriter, r *http.Request) {
    f, fHeader,err := r.FormFile("content");
    if err != nil {
        fmt.Fprintln(w, err);
        return;
    }
    defer f.Close();
    size := 0;
    download, err := os.OpenFile("./gokDownload", os.O_RDWR, 0644);
    if err != nil {
        fmt.Fprintln(w, err);
        return;
    }
    defer os.Close(download);
    buff := make([]byte, 512);
    for {
        n, err := f.Read(buff);
        if (n == 0) || (err == io.EOF) {
            break;
        }
        if err != nil {
            fmt.Fprintln(w, err);
            return;
        }
        n2, err := os.Write(download, buff);
        if err != nil {
            fmt.Fprintln(w, err);
        }
        size += n2;
    }
    fmt.Fprintln(w, "File Recived Successfully PATH =>", f2, " Size =>",size);
}

type mainHandler struct{};
func (_ *mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
   if r.URL.Path == "/" {
       index(w, r);
   } else if r.URL.Path == "/postit" {
       postit(w, r);
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
