package main
import "fmt"
import "net/http"
import "net/url"
import "time"
import "strings"
import "os"
import "io"
import "bytes"

// core struct encapsulates the http.Request and http.ResponseWriter. 
// instance present in all <?go ?> tags as `gok`
type Gok struct {
    w http.ResponseWriter
    r *http.Request
    getValues url.Values
    response *bytes.Buffer
    deadMsg string
    shouldRedirect bool
};

// equivalent to php `echo`. echo accepts any type of parameters
func (self *Gok) Echo(a ...interface{}) {
    if self.response != nil {
        fmt.Fprint(self.response, a...)
    }
}

// redirects the request to a new url 
func (self *Gok) Redirect(newUrl string) {
    self.shouldRedirect = true
    self.Header("Location:"+newUrl)
    self.w.WriteHeader(http.StatusMovedPermanently)
}

// Die sends the `msg` as error, and undoes everything echoed so far.
func (self *Gok) Die(msg string) {
    self.response.Reset()
    self.response = nil
    self.deadMsg = msg
}

//=> PHP equivalent of $_SERVER[] functions.
// returns the current request url without the leading `/`
func (self *Gok) ServerSelf() string {
    return self.r.URL.Path[1:]
}
func (self *Gok) ServerHttpUserAgent() string {
    return self.r.Header.Get("User-Agent")
}
func (self *Gok) ServerHttpReferer() string {
    return self.r.Referer()
}
func (self *Gok) ServerHttps() bool {
    return self.r.URL.Scheme == "https"
}
func (self *Gok) ServerRemoteAddr() string {
    return strings.Split(self.r.RemoteAddr, ":")[0]
}
func (self *Gok) ServerRemotePort() string {
    return strings.Split(self.r.RemoteAddr, ":")[1]
}
func (self *Gok) ServerPort() int {
    return 80
}
func (self *Gok) ServerHttpAcceptEncoding() string {
    return self.r.Header.Get("Accept-Encoding")
}
func (self *Gok) ServerProtocol() string {
    return self.r.Proto
}
func (self *Gok) ServerRequestMethod() string {
    return self.r.Method
}
func (self *Gok) ServerQueryString() string {
    return self.r.URL.RawQuery
}
func (self *Gok) ServerHttpAccept() string {
    return self.r.Header.Get("Accept")
}
func (self *Gok) ServerHttpAcceptCharset() string {
    return self.r.Header.Get("Accept-Charset")
}
func (self *Gok) ServerHttpAcceptLanguage() string {
    return self.r.Header.Get("Accept-Language")
}
func (self *Gok) ServerHttpConnection() string {
    return self.r.Header.Get("Connection")
}
func (self *Gok) ServerHttpHost() string {
    return self.r.Host
}


//=> PHP equivalent of $_GET[''] and $_POST['']
// parses the post request and returns the post value of `name`,
// in case for post containing file date. use gok.File('fileName')
func (self *Gok) Post(name string) string {
    return self.r.PostFormValue(name)
}
// the same as above but for get requests.
func (self *Gok) Get(name string) string {
    if self.getValues == nil {
        self.getValues = self.r.URL.Query()
    }
    return self.getValues.Get(name)
}

//=> PHP equivalent of $_COOKIE['']
// returns the cookie value for `name`, if no cookie is set returns an
// empty string ""
func (self *Gok) Cookie(name string) string {
    if cookie, err := self.r.Cookie(name); err == http.ErrNoCookie {
        return ""
    } else if cookie != nil {
        return cookie.Value
    }
    return ""
}
// sets cookie with name, value, and duration to expire, if duration is 0
// cookie is permanant
func (self *Gok) SetCookie(name string, value string, duration int64) {
    if (len(name) == 0) || (len(value) == 0) {
        return
    }
    cookie := new(http.Cookie)
    cookie.Name = name
    cookie.Value = value
    if duration != 0 {
        cookie.Expires = time.Now().Add(time.Duration(duration) * time.Second)
    }
    http.SetCookie(self.w, cookie)
}
//delets the cookie set with `name`
func (self *Gok) DeleteCookie(name string) {
    self.SetCookie(name, "deleted", -1000)
}

func (self *Gok) SetCookie_4(name string, value string, duration int64,
                                urlPath string){
    if (len(name) == 0) || (len(value) == 0) {
        return
    }
    cookie := new(http.Cookie)
    cookie.Name = name
    cookie.Value = value
    if duration != 0 {
        cookie.Expires = time.Now().Add(time.Duration(duration) * time.Second)
    }
    if len(urlPath) != 0 {
        cookie.Path = urlPath
    }
    http.SetCookie(self.w, cookie)
}

func (self *Gok) SetCookie_5(name string, value string, duration int64,
                                urlPath string, domain string) {
    if (len(name) == 0) || (len(value) == 0) {
        return
    }
    cookie := new(http.Cookie)
    cookie.Name = name
    cookie.Value = value
    if duration != 0 {
        cookie.Expires = time.Now().Add(time.Duration(duration) * time.Second)
    }
    if len(urlPath) != 0 {
        cookie.Path = urlPath
    }
    if len(domain) != 0 {
        cookie.Domain = domain
    }
    http.SetCookie(self.w, cookie)
}

func (self *Gok) SetCookie_7(name string, value string, duration int64,
                                urlPath string, domain string, secure bool,
                                httpOnly bool) {
    if (len(name) == 0) || (len(value) == 0) {
        return
    }
    cookie := new(http.Cookie)
    cookie.Name = name
    cookie.Value = value
    if duration != 0 {
        cookie.Expires = time.Now().Add(time.Duration(duration) * time.Second)
    }
    if len(urlPath) != 0 {
        cookie.Path = urlPath
    }
    if len(domain) != 0{
        cookie.Domain = domain
    }
    cookie.Secure = secure
    cookie.HttpOnly = httpOnly
    http.SetCookie(self.w, cookie)
}

//=> PHP equivalent of $FILE['']
// File saves the file upload to disk and returns 
// (`file path`,`uploaded file name`, `file content type`, `size`) of the file.
// its important that file uploads occure from form with 
// enctype='multipart/form-data' set.
func (self *Gok) File(name string) (string, string, string, int64) {
    f, fHeader, err := self.r.FormFile(name)
    if err != nil {
        return "", "", "", 0
    }
    defer f.Close()
    fileName := genRandName()
    f2, err := os.Create(fileName)
    if err != nil {
        return "", "", "", 0
    }
    defer f2.Close()
    size, err := io.Copy(f2, f)
    if err != nil {
        return "", "", "", 0
    }
    if len(fHeader.Header["Content-Type"]) == 0 {
        return fileName, fHeader.Filename, "", size
    }
    return fileName, fHeader.Filename, fHeader.Header["Content-Type"][0], size
}

//=> PHP equivalent of $_HEADER['']
// sets the header for response.
// ex: gok.Header("Connection:Close")
func (self *Gok) Header(header string) {
    h := strings.Split(header, ":")
    if len(h) != 2 {
        panic("unknown header value")
    }
    self.w.Header().Set(h[0], h[1])
}
//go htt.Header for request
func (self *Gok) RequestHeader() http.Header {
    return self.r.Header
}
//go http.Header for Response
func (self *Gok) ResponseHeader() http.Header {
    return self.w.Header()
}

// htt.ResponseWriter for gok object
func (self *Gok) ResponseWriter() http.ResponseWriter { return self.w }
// *http.Request for gok object
func (self *Gok) HttpRequest() *http.Request { return self.r }
