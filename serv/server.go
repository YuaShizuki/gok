package main
import "fmt"
import "net/http"


var routes map[string]func(*Gok) = map[string]func(*Gok) {
//<gok routers>
    "/":index,
//</gok>
};

func index(gok *Gok) {
    gok.Echo("<html>");
    gok.Echo("ServerSelf: ", gok.ServerSelf(),"</br>");
    gok.Echo("ServerHttpUserAgent: ", gok.ServerHttpUserAgent(),"</br>");
    gok.Echo("SeverHttps: ", gok.ServerHttps(), "</br>");
    gok.Echo("ServerRemoteAddr: ", gok.ServerRemoteAddr(), "</br>");
    gok.Echo("ServerRemotePort: ",gok.ServerRemotePort(), "</br>");
    gok.Echo("ServerPort: ", gok.ServerPort(), "</br>");
    gok.Echo("ServerHttpAcceptEncoding: ", gok.ServerHttpAcceptEncoding(), "</br>");
    gok.Echo("ServerProtocol: ",gok.ServerProtocol(),"</br>");
    gok.Echo("ServerRequestMethod: ",gok.ServerRequestMethod(),"</br>");
    gok.Echo("ServerQueryString: ",gok.ServerQueryString(),"</br>");
    gok.Echo("ServerHttpAccept: ",gok.ServerHttpAccept(),"</br>");
    gok.Echo("ServerHttpAcceptCharset: ",gok.ServerHttpAcceptCharset(),"</br>");
    gok.Echo("ServerHttpAcceptLanguage: ",gok.ServerHttpAcceptLanguage(),"</br>");
    gok.Echo("ServerHttpConnection: ",gok.ServerHttpConnection(),"</br>");
    gok.Echo("ServerHttpHost: ",gok.ServerHttpHost(),"</br>");
    gok.Echo("</html>");
}

func submit(gok *Gok) {
    gok.Echo(gok.File("content"));
}

type mainHandler struct{};
func (_ *mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    gok := new(Gok);
    gok.r = r;
    gok.w = w;
    fn, ok := routes[r.URL.Path];
    if ok {
        fn(gok);
        fmt.Fprint(w, gok.response);
    } else {
        w.WriteHeader(404);
        fmt.Fprintln(w, "Error Not Found!");
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
