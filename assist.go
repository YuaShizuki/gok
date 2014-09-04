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

func unpackResource(created *list.List) {
    orignal, err := hex.DecodeString(Resource)
    if err != nil {
        errExit(err, "")
    }
    unzip, err := gzip.NewReader(bytes.NewBuffer([]byte(orignal)))
    if err != nil {
        errExit(err, "")
    }
    untar := tar.NewReader(unzip)
    var fileContent bytes.Buffer
    for h, err := untar.Next(); err == nil; h, err = untar.Next() {
        fileContent.Reset()
        io.Copy(&fileContent, untar)
        var name string
        if h.Name[0] == '/' {
            name = h.Name[1:]
        } else {
            name = h.Name
        }
        if strings.Contains(name, "/") {
            lstIndx := strings.LastIndex(name, "/")
            if err := os.MkdirAll(name[0:lstIndx], 0744); err != nil {
                errExit(err, "")
            }
        }
        ioutil.WriteFile(name, fileContent.Bytes(), os.FileMode(h.Mode))
        if name[len(name)-3:] == ".go" {
            created.PushBack(name)
        }
    }
}
