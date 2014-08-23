package main
import "fmt"
import "os"
import "encoding/hex"
import "math/rand"
import "time"

func pathExist(p string) bool {
    _, err := os.Stat(p)
    return err == nil;
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

