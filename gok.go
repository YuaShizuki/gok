package main
import "fmt"
import "io/ioutil"
import "path/filepath"

var (
    fileExtension = "gok"
);

func main() {
    files,_:= filepath.Glob("./*."+fileExtension);
    for _,v := range files {
        content, err := ioutil.ReadFile(v);
        if err != nil {
            errExit(err, "");
        }
        go_code,err := processGok(string(content));
        if err != nil {
            errExit(err, "");
        }
        fmt.Println(go_code);
    }
}
