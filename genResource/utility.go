package main
import "fmt"
import "os"
import "encoding/hex"
import "math/rand"
import "time"


pathExist(path string) (bool, os.FileInfo) {
    info, err :=  os.Stat(path);
    if err != nil {
        if os.IsNotExist(err) {
            return false, info;
        }
        panic("panic-> cannot determine if path exists");
    }
    return true, info;
}

func errExit(e error, msg string) {
    fmt.Println("error =>", e.Error());
    if msg != "" {
        fmt.Println(msg);
    }
    os.Exit(1);
}

func genRandName() string {
    rand.Seed(time.Now().UnixNano());
    num := rand.Uint32();
    str := []byte{ byte(num), byte(num >> 8), byte(num >> 16), byte(num >> 24) };
    return hex.EncodeToString(str);
}

