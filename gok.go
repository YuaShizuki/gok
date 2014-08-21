package main
import "fmt"
import "os"
import "io/ioutil"
import "path/filepath"

var (
    fileExtension = "gok"
    createdFiles map[string]string;
);

func main() {
    //TODO remove this after stabilization. 
    if len(os.Args) == 2 {
        os.Chdir(os.Args[1]);
    }
    createdFiles = make(map[string]string);
    convertGokToGoFiles(".");
    for k,v := range createdFiles {
        fmt.Println(k,":",v);
    }
}

func convertGokToGoFiles(dir string) {
    files, err := filepath.Glob(dir+"/*.gok");
    if err != nil {
        errExit(err, "");
    }
    for _, s := range files {
        gokContent, err := ioutil.ReadFile(s);
        if err != nil {
            errExit(err, "");
        }
        gocode, mainFunc, err := processGokContent(string(gokContent));
        if err != nil {
            errExit(err, "");
        }
        createdFiles[s] = mainFunc;
        ioutil.WriteFile(s+".go", []byte(gocode), 0644);
    }
    if dirs, err := ioutil.ReadDir(dir); err == nil {
        for _, d := range dirs {
            if d.IsDir() {
                convertGokToGoFiles(dir+"/"+d.Name());
            }
        }
    }
}

