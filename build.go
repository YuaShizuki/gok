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
import "errors"
import "strconv"

var (
    fileExtension = "gok"
    webRoutes     map[string]string
    shouldDelete  *list.List
)

func build() error {
    webRoutes = make(map[string]string)
    shouldDelete = list.New()
    convertGokToGoFiles(".")
    unpackResource()
    injectRoutes()
    out := goBuild()
    defer delFiles(shouldDelete)
    if len(out) != 0 {
        out = convertGoErrorsToGok(out);
        return errors.New(out)
    }
    return nil;
}

func convertGokToGoFiles(dir string) {
    files, err := filepath.Glob(dir + "/*.gok");
    if err != nil {
        errExit(err, "");
    }
    for _, s := range files {
        gokContent, err := ioutil.ReadFile(s)
        if err != nil {
            errExit(err, "")
        }
        gocode, mainFunc, err := processGokContent(string(gokContent))
        if err != nil {
            errExit(err, "")
        }
        webRoutes[s] = mainFunc
        goFile := buildGoFileName(s)
        shouldDelete.PushBack(goFile)
        ioutil.WriteFile(goFile, []byte(gocode), 0644)
    }
    if dirs, err := ioutil.ReadDir(dir); err == nil {
        for _, d := range dirs {
            if d.IsDir() {
                convertGokToGoFiles(dir + "/" + d.Name())
            }
        }
    }
}

func unpackResource() {
    orignal, err := hex.DecodeString(Resource)
    if err != nil {
        errExit(err, "")
    }
    unzip, err := gzip.NewReader(bytes.NewBuffer([]byte(orignal)))
    if err != nil {
        errExit(err, "")
    }
    untar := tar.NewReader(unzip)
    var fileContent bytes.Buffer
    for h, err := untar.Next(); err == nil; h, err = untar.Next() {
        fileContent.Reset()
        io.Copy(&fileContent, untar)
        var name string
        if h.Name[0] == '/' {
            name = h.Name[1:]
        } else {
            name = h.Name
        }
        if strings.Contains(name, "/") {
            lstIndx := strings.LastIndex(name, "/")
            if err := os.MkdirAll(name[0:lstIndx], 0744); err != nil {
                errExit(err, "")
            }
        }
        ioutil.WriteFile(name, fileContent.Bytes(), os.FileMode(h.Mode))
        if name[len(name)-3:] == ".go" {
            shouldDelete.PushBack(name)
        }
    }
}

func injectRoutes() {
    var routes bytes.Buffer
    s, err := ioutil.ReadFile("serverb596f256.go")
    if err != nil {
        errExit(err, "")
    }
    parts := strings.Split(string(s), "//<gok inject routes>")
    if len(parts) != 2 {
        errExit(nil, "Unusual resource file")
    }
    for k, v := range webRoutes {
        if k == "index.gok" {
            routes.Write([]byte("\"\":" + v + ",\n"))
        }
        routes.Write([]byte("\"" + k + "\":" + v + ",\n"))
    }
    final := strings.Join([]string{parts[0], string(routes.Bytes()), parts[1]}, "\n")
    ioutil.WriteFile("serverb596f256.go", []byte(final), 0644)
}

func goBuild() string {
    cmd := exec.Command("go", "build")
    output, _ := cmd.CombinedOutput()
    if len(strings.TrimSpace(string(output))) == 0 {
        return ""
    }
    return string(output)
}

func convertGoErrorsToGok(output string) string {
    out := new(bytes.Buffer)
    lines := strings.Split(output, "\n")
    for _, l := range lines {
        lineNum := 0
        var err error
        if len(l) < 2 {
            continue
        }
        if l[0] == '#' {
            fmt.Fprintln(out, l)
            continue
        }
        indx := strings.Index(l, ":")
        if (indx == -1) || (len(l) <= (indx + 1)) {
            fmt.Fprintln(out, l)
            continue
        }
        file := l[0:indx]
        gokFile := getEquivalentGokFile(file)
        if gokFile == "" {
            fmt.Fprintln(out, l)
            continue
        }
        indx2 := strings.Index(l[indx+1:], ":") + indx + 1
        if (indx2 == -1) || (len(l) <= (indx2 + 1)) {
            fmt.Fprintln(out, l)
            continue
        }
        if lineNum, err = strconv.Atoi(l[indx+1 : indx2]); err != nil {
            fmt.Fprintln(out, l)
            continue
        }
        gokFileLn, err := gokFileLineNum(gokFile, file, lineNum)
        if err != nil {
            fmt.Fprintln(out, l)
            continue
        }
        fmt.Fprintf(out, "%s:%d:%s\n", gokFile, gokFileLn, l[indx2+1:])
    }
    return out.String();
}

func getEquivalentGokFile(file string) string {
    if !pathExist(file) {
        return ""
    }
    if indx := strings.Index(file, ".gok.go"); (len(file) - 7) != indx {
        return ""
    }
    gokFileTemp := file[:len(file)-3];
    gokFile := strings.Replace(gokFileTemp, "[", "/", -1);
    if !pathExist(gokFile) {
        return ""
    }
    return gokFile
}


func gokFileLineNum(gokFile string, file string, ln int) (int, error) {
    f, err := ioutil.ReadFile(file)
    if err != nil {
        return 0, err
    }
    gokFileContent, err := ioutil.ReadFile(gokFile)
    if err != nil {
        return 0, err
    }
    lines := strings.Split(string(f), "\n")
    if len(lines) < ln {
        return 0, errors.New("line numbers exceds existing lines")
    }
    duplicates := countDuplicatesTill(lines, ln)
    orignalLn := findLnInGokFile(string(gokFileContent), lines[ln-1], duplicates)
    return orignalLn, nil
}

func countDuplicatesTill(code []string, ln int) int {
    initial := code[ln-1]
    count := 0
    for _, c := range code {
        if count == ln-1 {
            return count
        }
        if c == initial {
            count++
        }
    }
    return count
}

func findLnInGokFile(gokContent string, ln string, skipCount int) int {
    last := 0
    for i := 0; i < (skipCount + 1); i++ {
        slice := gokContent[last:]
        indx := strings.Index(slice, ln)
        if indx == -1 {
            break
        }
        last += indx
    }
    if last == 0 {
        return 0
    }
    return strings.Count(gokContent[0:last], "\n") + 1
}

func buildGoFileName(p string) string {
    ret := strings.Replace(p, "/", "[", -1)
    return ret+".go"
}
