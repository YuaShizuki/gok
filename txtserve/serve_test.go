package txtserve
import "testing"
import "time"
import "net/http"
import "io/ioutil"
import "strings"

func getForTest() (string, error) {
    resp, err := http.Get("http://127.0.0.1/")
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(body)), nil
}


func TestServer(t *testing.T) {
    err := StartServer("TEST")
    if err != nil {
        t.Fatal(err)
    }
    time.Sleep(5 * time.Second)
    resp, err := getForTest()
    if (err != nil) || (resp != "TEST") {
        t.Fatal(err, resp)
    }
    err = StopServer()
    if err != nil {
        t.Fatal(err)
    }
    time.Sleep(5 * time.Second)
    resp, err = getForTest()
    if (err == nil)  || (resp == "TEST") {
        t.Fatal("server still running")
    }
}
