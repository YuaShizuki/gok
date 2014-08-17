package main
import "fmt"
import "net/http"
import "io"
import "os"
import "math/rand"
import "encoding/hex"
import "time"
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
    if _, err := r.Cookie("saltRocket"); err == http.ErrNoCookie {
        fmt.Println("No Cookies Recived");
        /*cookie := &http.Cookie{
                        Name:"saltRocket",
                        Value:"fossdark",
                    };*/
        fmt.Println("=> recived cookie");
    }
    fmt.Fprintln(w, html);
}

func clearCookie(w http.ResponseWriter, r *http.Request) {
    if cookie, err := r.Cookie("saltRocket"); err == http.ErrNoCookie {
        fmt.Println("No cookie present")
    } else if cookie != nil {
        cookie.Expires = time.Now().Add(-(300 * 60 * time.Minute));
        cookie.RawExpires = "";
        http.SetCookie(w, cookie);
    }
    fmt.Fprintln(w, "<html><h1>fuck me!.</h1></html>")
}

func genRandName() string {
    rand.Seed(time.Now().UnixNano());
    num := rand.Uint32();
    str := []byte{ byte(num), byte(num >> 8), byte(num >> 16), byte(num >> 24) };
    return hex.EncodeToString(str);
}

func postit(w http.ResponseWriter, r *http.Request) {
    f, _, err := r.FormFile("content");
    if err != nil {
        fmt.Fprintln(w, err);
        return;
    }
    defer f.Close();
    name := genRandName();
    f2, err := os.Create(name);
    if err != nil {
        fmt.Fprintln(w, err);
        return;
    }
    defer f2.Close();
    io.Copy(f2, f);
    fmt.Fprintln(w, name);
}

type mainHandler struct{};
func (_ *mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
   if r.URL.Path == "/" {
       index(w, r);
   } else if r.URL.Path == "/postit" {
       postit(w, r);
   } else if r.URL.Path == "/clearcookie" {
        clearCookie(w, r);
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
