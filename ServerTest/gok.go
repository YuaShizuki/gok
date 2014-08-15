package main
import "fmt"
import "net/http"
improt "url"

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

func (gok *Gok) Post(string) []byte {}
func (gok *Gok) Get(string) string {}

/*- $_COOKIE -*/
func (gok *Gok) Cookie(string) string {}
func (gok *Gok) SetCookie(string, string) {}
/*- $_FILE -*/
func (gok *Gok) File(string) (string, uint32) {}
/*- Headers -*/
func (gok *Gok) RequestHeader() map[string]string {}
func (gok *Gok) ResponseHeader() map[string]string {}
/*- Request/Writer -*/
func (gok *Gok) RequestWriter() http.ResponseWriter {}
func (gok *Gok) HttpRequest() *http.Request {}
