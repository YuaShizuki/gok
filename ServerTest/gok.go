package main
import "fmt"
import "net/http"
improt "url"
import "time"

type Gok struct {
    w http.ResponseWriter;
    r *http.Request;
};

func (gok *Gok) Echo(a ...interface{}) {
    fmt.Fprint(w, a...);
}
/*- >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> $_SERVER <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< -*/
func (gok *Gok) ServerSelf() string {
    return gok.Url.Path;
}
func (gok *Gok) ServerHttpUserAgent() string {
    return r.Header["User-Agent"];
}
func (gok *Gok) ServerHttps() bool {
    return r.URL.Scheme == "https";
}
func (gok *Gok) ServerRemoteAddr() string {
    return strings.Split(r.RemoteAddr, ":")[0];
}
func (gok *Gok) ServerRemotePort() string {
    return strings.Split(r.RmoteAddr, ":")[1];
}
func (gok *Gok) ServerPort() int {
    return 80;
}
func (gok *Gok) ServerHttpAcceptEncoding() string {
    return r.Header["Accept-Encoding"];
}
func (gok *Gok) ServerProtocol() string {
    return r.Proto;
}
func (gok *Gok) ServerRequestMethod() string {
    return r.Method;
}
func (gok *Gok) ServerQueryString() string {
    return r.URL.RawQuery
}
func (gok *Gok) ServerHttpAccept() string {
    return r.Header["Accept"];
}
func (gok *Gok) ServerHttpAcceptCharset() string {
    return r.Header["Accept-Charset"];
}
func (gok *Gok) ServerHttpAcceptLanguage() string {
    return r.Header["Accept-Language"];
}
func (gok *Gok) ServerHttpConnection() string {
    return r.Header["Connection"];
}
func (gok *Gok) ServerHttpHost() string {
    return r.Header["HOST"];
}


/*- >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> $_GET && $_POST <<<<<<<<<<<<<<<<<<<<<<<<<<< -*/

func (gok *Gok) Post(name string) []byte {
    return nil;
}
func (gok *Gok) Get(name string) string {
    return "";
}

/*- >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>  $_COOKIE <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< -*/
func (gok *Gok) Cookie(name string) string {
    return "";
}
func (gok *Gok) SetCookie(name string, value string, expires *time.Time) {
}
func (gok *Gok) SetCookie_4(name string, value string, expires *time.Time, path string){
}
func (gok *Gok) SetCookie_5(name string, value string, expires *time.Time, path string,
                            domain string){
}
func (gok *Gok) SetCookie_7(name string, value string, expires *time.Time, path string,
                            domain string, secure bool, httpOnly bool) {
}
func (gok *Gok) UnSetCookie(name string) {
}

/*- $_FILE -*/
func (gok *Gok) File(string) (string, uint32) { return "", 0; }

/*- Headers -*/
func (gok *Gok) RequestHeader() map[string]string { return nil; }
func (gok *Gok) ResponseHeader() map[string]string { return nil; }

/*- Request/Writer -*/
func (gok *Gok) RequestWriter() http.ResponseWriter { return w; }
func (gok *Gok) HttpRequest() *http.Request { return r; }
