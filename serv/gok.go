package main
import "fmt"
import "net/http"
import "time"
import "strings"
import "os"
import "io"

type Gok struct {
    w http.ResponseWriter;
    r *http.Request;
};

func (self *Gok) Echo(a ...interface{}) {
    fmt.Fprint(self.w, a...);
}
/*- >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> $_SERVER <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< -*/
func (self *Gok) ServerSelf() string {
    return self.r.URL.Path;
}
func (self *Gok) ServerHttpUserAgent() string {
    return strings.Join(self.r.Header["User-Agent"], "\n");
}
func (self *Gok) ServerHttpReferer() string {
    return strings.Join(self.r.Header["Referer"], "\n");
}

func (self *Gok) ServerHttps() bool {
    return self.r.URL.Scheme == "https";
}
func (self *Gok) ServerRemoteAddr() string {
    return strings.Split(self.r.RemoteAddr, ":")[0];
}
func (self *Gok) ServerRemotePort() string {
    return strings.Split(self.r.RemoteAddr, ":")[1];
}
func (self *Gok) ServerPort() int {
    return 80;
}
func (self *Gok) ServerHttpAcceptEncoding() string {
    return strings.Join(self.r.Header["Accept-Encoding"], "\n");
}
func (self *Gok) ServerProtocol() string {
    return self.r.Proto;
}
func (self *Gok) ServerRequestMethod() string {
    return self.r.Method;
}
func (self *Gok) ServerQueryString() string {
    return self.r.URL.RawQuery
}
func (self *Gok) ServerHttpAccept() string {
    return strings.Join(self.r.Header["Accept"], "\n");
}
func (self *Gok) ServerHttpAcceptCharset() string {
    return strings.Join(self.r.Header["Accept-Charset"], "\n");
}
func (self *Gok) ServerHttpAcceptLanguage() string {
    return strings.Join(self.r.Header["Accept-Language"], "\n");
}
func (self *Gok) ServerHttpConnection() string {
    return strings.Join(self.r.Header["Connection"], "\n");
}
func (self *Gok) ServerHttpHost() string {
    return self.r.Host;
}


/*- >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> $_GET && $_POST <<<<<<<<<<<<<<<<<<<<<<<<<<< -*/

func (self *Gok) Post(name string) []byte {
    return nil;
}
func (self *Gok) Get(name string) string {
    return "";
}

/*- >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>  $_COOKIE <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< -*/
func (self *Gok) Cookie(name string) string {
    return "";
}
func (self *Gok) SetCookie(name string, value string, expires *time.Time) {
}
func (self *Gok) SetCookie_4(name string, value string, expires *time.Time, path string){
}
func (self *Gok) SetCookie_5(name string, value string, expires *time.Time, path string,
                            domain string){
}
func (self *Gok) SetCookie_7(name string, value string, expires *time.Time, path string,
                            domain string, secure bool, httpOnly bool) {
}
func (self *Gok) UnSetCookie(name string) {
}

/*- $_FILE -*/
func (self *Gok) File(name string) (string, string, string, int64) {
    f, fHeader, err := self.r.FormFile(name);
    if err != nil {
        return "", "", "", 0;
    }
    defer f.Close();
    fileName := genRandName();
    f2, err := os.Create(fileName);
    if err != nil {
        return "", "", "", 0;
    }
    defer f2.Close();
    size, err := io.Copy(f2, f);
    if err != nil {
        return "", "", "", 0;
    }
    if len(fHeader.Header["Content-Type"]) == 0 {
        return fileName, fHeader.Filename, "", size;
    }
    return fileName, fHeader.Filename, fHeader.Header["Content-Type"][0], size;
}

/*- Headers -*/
func (self *Gok) Header(header string) {
    h := strings.Split(header, ":");
    if len(h) != 2 {
        panic("unknown header value");
    }
    self.w.Header().Add(h[0], h[1]);
}

func (self *Gok) RequestHeader() map[string][]string {
    return self.r.Header;
}
func (self *Gok) ResponseHeader() map[string][]string {
    return self.w.Header();
}

/*- Request | Writer -*/
func (self *Gok) ResponseWriter() http.ResponseWriter { return self.w; }
func (self *Gok) HttpRequest() *http.Request { return self.r; }
