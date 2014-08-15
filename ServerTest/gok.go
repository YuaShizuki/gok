package main
import "fmt"
import "net/http"
import "time"
import "strings"

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
    return strings.Join(self.r.Header["User-Agent"], " ");
}
func (self *Gok) ServerHttpReferer() string {
    return strings.Join(self.r.Header["Referer"], " ");
}

func (self *Gok) ServerHttps() bool {
    return self.r.URL.Scheme == "https";
}
func (self *Gok) ServerRemoteAddr() string {
    return strings.Split(self.r.RemoteAddr, ":")[0];
}
func (self *Gok) ServerRemotePort() string {
    return strings.Split(self.r.RmoteAddr, ":")[1];
}
func (self *Gok) ServerPort() int {
    return 80;
}
func (self *Gok) ServerHttpAcceptEncoding() string {
    return strings.Join(self.r.Header["Accept-Encoding"], " ");
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
    return strings.Join(self.r.Header["Accept"], " ");
}
func (self *Gok) ServerHttpAcceptCharset() string {
    return strings.Join(self.r.Header["Accept-Charset"], " ");
}
func (self *Gok) ServerHttpAcceptLanguage() string {
    return strings.Join(self.r.Header["Accept-Language"], " ");
}
func (self *Gok) ServerHttpConnection() string {
    return self.r.Header["Connection"];
}
func (self *Gok) ServerHttpHost() string {
    return self.r.Header["Host"];
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
func (self *Gok) File(string) (string, uint32) { return "", 0; }

/*- Headers -*/
func (self *Gok) RequestHeader() map[string]string { return nil; }
func (self *Gok) ResponseHeader() map[string]string { return nil; }

/*- Request/Writer -*/
func (self *Gok) RequestWriter() http.ResponseWriter { return w; }
func (self *Gok) HttpRequest() *http.Request { return r; }
