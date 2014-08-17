package main
import "fmt"
import "net/http"
import "net/url"
import "time"
import "strings"
import "os"
import "io"
import "bytes"

type Gok struct {
    w http.ResponseWriter;
    r *http.Request;
    getValues url.Values;
    response *bytes.Buffer;
    should_redirect bool;
};

func (self *Gok) Echo(a ...interface{}) {
    if response = nil {
        self.response = new(bytes.Buffer);
    }
    fmt.Fprint(self.response, a...);
}

func (self *Gok) Redirect(newUrl string) {
    self.should_redirect = true;
    self.Header("Location:"+url);
}

func (self *Gok) Die() {
    self.response = nil;
}

/*- >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> $_SERVER <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< -*/
func (self *Gok) ServerSelf() string {
    return self.r.URL.Path;
}
func (self *Gok) ServerHttpUserAgent() string {
    return strings.Join(self.r.Header["User-Agent"], "\n");
}
func (self *Gok) ServerHttpReferer() string {
    self.r.Referer();
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
    return self.r.PostFormValue(name);
}
func (self *Gok) Get(name string) string {
    if self.getValues == nil {
        self.getValues, err = url.ParseQuery(r.URL.Query);
        if err == nil {
            self.getValues = nil;
            return "";
        }
    }
    return self.getValues.Get(name);
}

/*- >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>  $_COOKIE <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< -*/
func (self *Gok) Cookie(name string) string {
    if cookie, err := self.r.Cookie(name); err == http.ErrNoCookie {
        return "";
    } else if cookie != nil {
        return cookie.Value;
    }
    return "";
}
func (self *Gok) SetCookie(name string, value string, duration int64) {
    if (len(name) == 0) || (len(value) == 0) {
        return;
    }
    cookie := new(http.Cookie);
    cookie.Name = name;
    cookie.Value = value;
    if duration != 0 {
        cookie.Expires = time.Now().Add(duration * time.Second);
    }
    http.SetCookie(w, cookie);
}

func (self *Gok) SetCookie_4(name string, value string, duration int64,
                                urlPath string){
    if (len(name) == 0) || (len(value) == 0) {
        return;
    }
    cookie := new(http.Cookie);
    cookie.Name = name;
    cookie.Value = value;
    if duration != 0 {
        cookie.Expires = time.Now().Add(duration * time.Second);
    }
    if len(urlPath) {
        cookie.Path = urlPath;
    }
    http.SetCookie(w, cookie);
}
func (self *Gok) SetCookie_5(name string, value string, duration int64,
                                urlPath string, domain string) {
    if (len(name) == 0) || (len(value) == 0) {
        return;
    }
    cookie := new(http.Cookie);
    cookie.Name = name;
    cookie.Value = value;
    if duration != 0 {
        cookie.Expires = time.Now().Add(duration * time.Second);
    }
    if len(urlPath) {
        cookie.Path = urlPath;
    }
    if len(domain) {
        cookie.Domain = domain;
    }
    http.SetCookie(w, cookie);
}
func (self *Gok) SetCookie_7(name string, value string, duration int64,
                                urlPath string, domain string, secure bool,
                                httpOnly bool) {
    if (len(name) == 0) || (len(value) == 0) {
        return;
    }
    cookie := new(http.Cookie);
    cookie.Name = name;
    cookie.Value = value;
    if duration != 0 {
        cookie.Expires = time.Now().Add(duration * time.Second)
    }
    if len(urlPath) {
        cookie.Path = urlPath;
    }
    if len(domain) {
        cookie.Domain = domain;
    }
    cookie.Secure = secure;
    cookie.HttpOnly = httpOnly;
    http.SetCookie(w, cookie);
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

func (self *Gok) RequestHeader() http.Header {
    return self.r.Header;
}
func (self *Gok) ResponseHeader() http.Header {
    return self.w.Header();
}

/*- Request | Writer -*/
func (self *Gok) ResponseWriter() http.ResponseWriter { return self.w; }
func (self *Gok) HttpRequest() *http.Request { return self.r; }
