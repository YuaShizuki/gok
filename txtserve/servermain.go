package main
import "time"
import "net/http"
import "io/ioutil"
import "strings"
import "fmt"

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


func TestAll() {
    err := StartServer("TEST")
    if err != nil {
        fmt.Println(err)
        return
    }
    response, err := getForTest()
    if (err != nil) || (response != "TEST") {
        t.Fatal(err, response)
    }
    time.Sleep(5 * time.Second)
    err := StopServer()
    time.Sleep(5 * time.Second)
    response, err := getForTest()
    if (err == nil)  || (response != "") {
        t.Fatal("server still running")
    }
}
