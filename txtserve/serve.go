package txtserve
import "fmt"
import "net"
import "errors"
import "strings"

var running bool = false
var response string
var header string = "HTTP/1.1 200 OK\r\nContent-Type: text; charset=UTF-8\r\nContent-Length: %d\r\nConnection: Close\r\n\r\n"
var mainListener net.Listener

func serv() {
    for ;; {
        conn, err := mainListener.Accept()
        if (err != nil) || (conn == nil) {
            return
        }
        temp := make([]byte, 512)
        for ;; {
            n, err := conn.Read(temp)
            if (n == 0) || (err != nil) {
                break
            }
            resp := string(temp[0:n])
            if strings.Index(resp, "\r\n\r\n") != -1 {
                break
            }
        }
        _, err = conn.Write([]byte(fmt.Sprintf(header, len(response)+2)))
        if err != nil { continue }
        _, err = conn.Write([]byte(response))
        if err != nil { continue }
        _, err = conn.Write([]byte("\r\n"))
        if err != nil { continue }
        conn.Close()
    }
}


func StartServer(content string) error {
    l, err := net.Listen("tcp", ":80")
    if err != nil {
        return err
    }
    mainListener = l
    response = content
    running = true
    go serv()
    return nil
}

func StopServer() error {
    running = false
    if mainListener == nil {
        return errors.New("server not listening")
    }
    err := mainListener.Close()
    return err
}

func IsRunning() bool {
    return running
}
