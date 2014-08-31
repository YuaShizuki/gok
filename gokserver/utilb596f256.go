package main
import "time"
import "math/rand"
import "encoding/hex"
import "os"
import "fmt"

func pathExist(path string) (bool, os.FileInfo) {
    info, err :=  os.Stat(path);
    if err != nil {
        return false, info;
    }
    return true, info;
}


func genRandName() string {
    rand.Seed(time.Now().UnixNano());
    num := rand.Uint32();
    str := []byte{ byte(num), byte(num >> 8), byte(num >> 16), byte(num >> 24) };
    return hex.EncodeToString(str);
}

func errExit(err error) {
    fmt.Println(err)
    os.Exit(1)
}
