package main
import "fmt"
import "io/ioutil"
import "path/filepath"

var (
    fileExtension = "gok"
);

func main() {
    files,_:= filepath.Glob("./*."+fileExtension);
    funcNames := make([]string, len(files));
    for i, v := range files {
        content, err := ioutil.ReadFile(v);
        if err != nil {
            errExit(err, "");
        }
        goCode, funcName ,err := processGok(string(content));
        if err != nil {
            errExit(err, "");
        }
        funcNames[i] = funcName;
        fmt.Println(goCode);
    }
    fmt.Println(funcNames)
}
