package main
import "time"
import "math/rand"
import "encoding/hex"

func genRandName() string {
    rand.Seed(time.Now().UnixNano());
    num := rand.Uint32();
    str := []byte{ byte(num), byte(num >> 8), byte(num >> 16), byte(num >> 24) };
    return hex.EncodeToString(str);
}
