package main
import "fmt"
import "os"
import "io"
import "io/ioutil"
import "path/filepath"
import "encoding/hex"
import "compress/gzip"
import "archive/tar"
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
    unpackResource();
    injectRoutes();
    fmt.Println("Gok exported golang code, execute go build")
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
    orignal, err := hex.DecodeString(Resource);
    if err != nil {
        errExit(err, "");
    }
    unzip, err := gzip.NewReader(bytes.NewBuffer([]byte(orignal)))
    if err != nil {
        errExit(err, "");
    }
    untar := tar.NewReader(unzip);
    var fileContent bytes.Buffer;
    for h, err := untar.Next(); err == nil;  h, err = untar.Next() {
        fileContent.Reset();
        io.Copy(&fileContent, untar);
        var name string;
        if h.Name[0] == '/' {
            name = h.Name[1:];
        } else {
            name = h.Name;
        }
        if strings.Contains(name, "/") {
            lstIndx := strings.LastIndex(name, "/");
            if err := os.MkdirAll(name[0:lstIndx], 0744); err != nil {
                errExit(err, "");
            }
        }
        ioutil.WriteFile(name, fileContent.Bytes(), os.FileMode(h.Mode));
    }
}

func injectRoutes() {
    var routes bytes.Buffer;
    s, err := ioutil.ReadFile("serverb596f256.go");
    if err != nil {
        errExit(err, "");
    }
    parts := strings.Split(string(s), "//<gok inject routes>");
    if len(parts) != 2 {
        errExit(nil, "Unusual resource file");
    }
    for k,v := range createdFiles {
        if k == "index.gok" {
            routes.Write([]byte("\"\":"+v+",\n"));
        }
        routes.Write([]byte("\""+k+"\":"+v+",\n"));
    }
    final := strings.Join([]string{parts[0], string(routes.Bytes()), parts[1]}, "\n");
    ioutil.WriteFile("serverb596f256.go", []byte(final), 0644);
}
