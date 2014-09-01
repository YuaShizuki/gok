package main
import "fmt"
import "os"
import "encoding/hex"
import "math/rand"
import "time"
import "container/list"
import "net"
import "io/ioutil"

func delFiles(l *list.List) {
    for e := l.Front(); e != nil; e = e.Next() {
        os.Remove(e.Value.(string))
    }
}

func pathExist(p string) bool {
    _, err := os.Stat(p)
    return err == nil
}

func errExit(e error, msg string) {
    fmt.Println("error =>", e.Error())
    if msg != "" {
        fmt.Println(msg)
    }
    os.Exit(1)
}

func genRandName() string {
    rand.Seed(time.Now().UnixNano())
    num := rand.Uint32()
    str := []byte{byte(num), byte(num >> 8), byte(num >> 16), byte(num >> 24)}
    return hex.EncodeToString(str)
}

func printUsage() {
    fmt.Println("#useage");
    fmt.Println("   $gok run");
    fmt.Println("   $gok build");
    fmt.Println("   $gok src");
    fmt.Println("   $gok api\n");
    os.Exit(1);
}

func waitTillPortShutDown(port string) {
    conn, err := net.Dial("tcp", "127.0.0.1"+port)
    if err != nil {
        return
    }
    defer conn.Close()
    ioutil.ReadAll(conn)
    return
}
