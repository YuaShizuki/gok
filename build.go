/*
- Copyright 2014 Yua Shizuki. All rights reserved.
- Licensed under the Apache License, Version 2.0 (the "License");
- you may not use this file except in compliance with the License.
- You may obtain a copy of the License at
-
- http://www.apache.org/licenses/LICENSE-2.0
-
- Unless required by applicable law or agreed to in writing, software
- distributed under the License is distributed on an "AS IS" BASIS,
- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
- See the License for the specific language governing permissions and
- limitations under the License.
-*/

package main
import "fmt"
import "io/ioutil"
import "path/filepath"
import "strings"
import "bytes"
import "os/exec"
import "container/list"
import "errors"
import "strconv"

var (
    fileExtension = "gok"
    webRoutes       map[string]string
    quickAjax       map[string]string
    shouldDelete    *list.List
)

func build(src bool) error {
    webRoutes = make(map[string]string)
    quickAjax = make(map[string]string)
    shouldDelete = list.New()
    err := convertGokToGoFiles(".")
    if err != nil {
        return err
    }
    unpackResource(shouldDelete)
    injectRoutes()
    if len(quickAjax) != 0 {
        buildGokJs(quickAjax)
        injectAjxRoutes(quickAjax)
    }
    if src {
        return nil
    }
    out := goBuild()
    defer delFiles(shouldDelete)
    if len(out) != 0 {
        out = convertGoErrorsToGok(out)
        return errors.New(out)
    }
    return nil;
}

func convertGokToGoFiles(dir string) error {
    files, err := filepath.Glob(dir + "/*.gok")
    if err != nil {
        errExit(err, "")
    }
    for _, s := range files {
        gokContent, err := ioutil.ReadFile(s)
        if err != nil {
            errExit(err, "")
        }
        gocode, mainFunc, ajxFuncs, err := compile(string(gokContent))
        if err != nil {
            return errors.New(s + " :"+err.Error())
        }
        appendToQuickAjax(ajxFuncs)
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
    return nil
}

func injectRoutes() {
    var routes bytes.Buffer
    s, err := ioutil.ReadFile("serve.auto.go")
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
    ioutil.WriteFile("serve.auto.go", []byte(final), 0644)
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
        gokFileLn := gokFileLineNum(file, lineNum)
        fmt.Fprintf(out, "%s:%s:%s\n", gokFile, gokFileLn, l[indx2+1:])
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


func gokFileLineNum(file string, ln int) string {
    f, err := ioutil.ReadFile(file)
    if err != nil {
        fmt.Println("gokFileLineNum Error => ",err)
        return "0"
    }
    lines := strings.Split(string(f), "\n")
    comment := lines[ln-2]
    if len(comment) < 2 {
        return "[unknown ln]"
    }
    return comment[2:]
}

func buildGoFileName(p string) string {
    ret := strings.Replace(p, "/", "[", -1)
    return ret+".go"
}

func appendToQuickAjax(newfuncs map[string]string) {
    for k,v := range newfuncs {
        quickAjax[k] = v
    }
}
