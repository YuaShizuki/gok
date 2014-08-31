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
    response, err := getForTest()
    if (err != nil) || (response != "TEST") {
        t.Fatal(err, response)
    }
    err = StopServer()
    if err != nil {
        t.Fatal(err)
    }
    time.Sleep(5 * time.Second)
    response, err = getForTest()
    if (err == nil)  || (response == "TEST") {
        t.Fatal("server still running")
    }
}
