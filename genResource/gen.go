package main
import "fmt"
import "io/ioutil"
import "bytes"
import "archive/tar"
import "compress/gzip"
import "encoding/hex"

func main() {
    if exist, _ := pathExist("resource.go"); exist {
        fmt.Println("resource.go already exists");
        return;
    }
    buff := new(bytes.Buffer);
    tw := tar.NewWriter(buff);
    err := packDir("", tw);
    if err != nil {
        errExit(err, "");
    }
    if err := tw.Close(); err != nil {
        errExit(err, "");
    }
    finalBuff := new(bytes.Buffer);
    gw := gzip.NewWriter(finalBuff);
    gw.Write(buff.Bytes());
    gw.Close();
    //build a hex string for writing to a .go file
    hexStr := hex.EncodeToString(finalBuff.Bytes());
    f, err := os.Create("resource.go");
    if err != nil {
        errExit(err);
    }
    f.Write([]byte("var Resource string = `"));
    f.Write([]byte(hexStr));
    f.Write([]byte("`;\n"));
    f.Close();
}

func packDir() {
}
