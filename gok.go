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
    for _,v := range files {
        content, err := ioutil.ReadFile(v);
        if err != nil {
            errExit(err, "");
        }
        goCode,_,err := processGok(string(content));
        if err != nil {
            errExit(err, "");
        }
        fmt.Println(goCode);
        ioutil.WriteFile(v+".go", []byte(goCode), 0644)
    }
    fmt.Println(funcNames)
}
