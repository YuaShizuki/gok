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
import "os/exec"
import "container/list"


var (
    fileExtension = "gok"
    webRoutes map[string]string;
    shouldDelete *list.List;
);

func main() {
    //TODO remove this after stabilization. 
    if len(os.Args) >= 2 {
        os.Chdir(os.Args[1]);
    }
    webRoutes = make(map[string]string);
    shouldDelete = list.New();
    convertGokToGoFiles(".");
    unpackResource();
    injectRoutes();
    fmt.Println("excuting go build");
    if len(os.Args) >= 3 {
        if os.Args[2] == "-testsrc" {
            return;
        }
    }
    out := goBuild();
    if len(out) > 0 {
        printFormatedBuildErrors(out);
        return;
    }
    deleteBuiltFiles();
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
        webRoutes[s] = mainFunc;
        shouldDelete.PushBack(s+".go");
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
        if name[len(name)-3:] == ".go" {
            shouldDelete.PushBack(name);
        }
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
    for k,v := range webRoutes {
        if k == "index.gok" {
            routes.Write([]byte("\"\":"+v+",\n"));
        }
        routes.Write([]byte("\""+k+"\":"+v+",\n"));
    }
    final := strings.Join([]string{parts[0], string(routes.Bytes()), parts[1]}, "\n");
    ioutil.WriteFile("serverb596f256.go", []byte(final), 0644);
}

func goBuild() {
    cmd := exec.Command("go", "build");
    output, _ := cmd.Output();
    if len(strings.TrimSpace(output)) == 0 {
        return "";
    }
    return output;
}

func deleteBuiltFiles() {
    for e := shouldDelete.Front(); e != nil; e = e.Next() {
       os.Remove(e.Value.(string));
   }
}

func printFormatedBuildErrors(output string) {
    out := new(bytes.Buffer);
    lines := strings.Split(output, "\n");
    for _, l := range lines {
        lineNum := 0;
        if len(l) < 2 { continue; }
        if l[0] == '#' {
            fmt.Fprintln(out, l);
            continue;
        }
        indx := strings.Index(l, ":");
        if (indx == -1) || (len(l) <= (indx+1)) {
            fmt.Fprintln(out, l);
            continue;
        }
        file := l[0:indx];
        gokFile := getEquivalentGokFile(file);
        if gokFile  == "" {
            fmt.Fprintln(out, l);
            continue;
        }
        indx2 := strings.Index(l[indx+1:], ":")+indx+1;
        if (indx2 == -1) || (len(l) <= (indx2 + 1)) {
            fmt.Fprintln(out, l);
            continue;
        }
        if lineNum, err := strconv.Atoi(l[indx+1:indx2]); err != nil {
            fmt.Fprintln(out, l);
            continue;
        }
        gokFileLn, err := gokFileLineNum(gokFile, file, lineNum);
        if err != nil {
            fmt.Fprintln(out, l);
            continue;
        }
        fmt.Fprintf(out, "%s:%d:%s\n", gokFile, gokFileLn, l[indx2+1:]);
    }
    return out.String();
}

func getEquivalentGokFile(file) bool {
    if !pathExist(file) {
        return "";
    }
    if indx := strings.Index(file, ".gok.go"); (len(v)-7) != indx {
        return "";
    }
    gokFile := file[:len(file)-3];
    if !pathExist(gokFile) {
        return "";
    }
    return gokFile;
}

func gokFileLineNum(file string, lineNum int) (int, error) {
    f, err := ioutil.ReadFile(file);
    if err != nil {
        return 0, err;
    }
    lines := strings.Split(string(f), "\n");
}

func countDuplicatesTill(code []string, ln int) {

}

func findLnInGokFile(gokContent string, ln string, skipCount int) int {

}
