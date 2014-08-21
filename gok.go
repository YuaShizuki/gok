package main
import "fmt"
import "os"
import "io/ioutil"
import "path/filepath"
import "encoding/hex"
import "compress/gzip"
import "archivei/tar"
import "strings"
import "bytes"

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
    unpackResource();
    injectRoutes();
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

func unpackResource() {
    orignal := hex.DecodeString(Resource);
    if err != nil {
        errExit(err, "");
    }
    unzip, err := gzip.NewReader(bytes.NewBuffer([]byte(orignal)))
    if err != nil {
        errExit(err);
    }
    untar := tar.NewReader(unzip);
    var fileContent bytes.Buffer;
    for h, err := untar.Next(); err == nil;  h, err = untar.Next() {
        fileContent.Reset();
        io.Copy(&fileContent, untar);
        var name string;
        if h[0] == '/' {
            name = h.Name[1:];
        } else {
            name = h.Name;
        }
        if strings.Contains(name, "/") {
            lstIndex := strings.LastIndex(name, "/");
            if err := os.MkdirAll(name[0:lstIndx], 0744); err != nil {
                errExit(err);
            }
        }
        ioutil.WriteFile(name, fileContent.Bytes(), os.FileMode(h.Mode));
    }
}

func injectRoutes() {
    s, err := ioutil.ReadFile("serverb596f256.go");
    if err != nil {
        errExit(err);
    }
    part1 := strings.Split(string(s), "//<gok routes>\n");
}
