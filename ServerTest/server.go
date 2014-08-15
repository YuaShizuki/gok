package main
import "fmt"
import "net/http"
import "io/ioutil"

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
    f,_,_ := r.FormFile("content");
    file,_ := ioutil.ReadAll(f);
    fmt.Fprintln(w, string(file));
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
