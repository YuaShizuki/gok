package main
import "fmt"
import "strings"
import "io/ioutil"
import "bytes"
import "errors"
import "container/list"

func main() {
    f, err := ioutil.ReadFile("i.gok")
    code, function, err := compile(string(f))
    if err != nil {
        fmt.Println(err)
        return;
    }
    fmt.Println(code)
    fmt.Println(function)
}

func compile(code string) (string, string, error) {
    lines := strings.Split(code, "\n")
    processed := list.New()
    header, err := extract("<?gohead", lines, takenIndxs)
    if err != nil { return "", "", err }
    funcs, err := extract("<?gofunc", lines, takenIndxs)
    if err != nil { return "", "", nil }
    code, err := compileRenderer(lines, takenIndxs)
    if err != nil { return "", "", nil }
    return header, "salvadorDali", err
}

type analyzed struct {
    ln int
    col1 int
    col2 int
}

func containsln(processed *list.List, ln int) (col1, col2 int) {
    for e := processed.Front; e != nil; e = e.Next() {
        lnsAndCols := e.Value.(*analyzed)
        if lnsAndCols.ln == ln {
            return lnsAndCols.col1, lnsAndCols.col2
        }
    }
    return -1, -1
}

func extract(pattern string, lines []string, processed *list.List) (string, error) {
    gocode := new(bytes.Buffer)
    plen := len(pattern)
    inside := false
    for i, l := range lines {
        lenl := len(l)
        if (lenl == 0) || (strings.TrimSpace(l) == "") {
            continue
        }
        if !inside {
            indx := strings.Index(l, pattern)
            if indx == -1 {
                continue
            }
            inside = true
            lnsAndCols := new(analyzed)
            lnsAndCols.ln = i
            if (lenl == plen) || (strings.TrimSpace(l[indx+plen:]) == "") {
                lnsAndCols.col1 = 0
                lnsAndCols.col2 = lenl
                processed.PushBack(lnsAndCols)
                continue
            }
            fmt.Fprintf(gocode, "//%d\n%s\n",(i+1), strings.TrimSpace(l[indx+plen:]))
            lnsAndCols.col1 = indx+plen
            lnsAndCols.col2 = lenl
            processed.PushBack(lnsAndCols)
        } else {
            lnsAndCols := new(analyzed)
            lnsAndCols.ln = i
            indx := strings.Index(l, "?>")
            if indx == -1 {
                fmt.Fprintf(gocode, "//%d\n%s\n", (i+1), strings.TrimSpace(l))
                lnsAndCols.col1 = 0
                lnsAndCols.col2 = lenl
                processed.PushBack(lnsAndCols)
                continue
            }
            inside = false;
            if lenl == 2 {
                lnsAndCols.col1 = 0
                lnsAndCols.col2 = lenl
                processed.PushBack(lnsAndCols)
                continue
            }
            if (indx == 0) || (strings.TrimSpace(l[:indx]) == "") {
                lnsAndCols.col1 = 0
                lnsAndCols.col1 = lenl
                processed.PushBack(lnsAndCols)
                continue
            }
            fmt.Fprintf(gocode, "//%d\n%s\n", (i+1), strings.TrimSpace(l[:indx]))
            lnsAndCols.col1 = 0
            lnsAndCols.col2 = indx
            processed.PushBack(lnsAndCols)
        }
    }
    if inside {
        return "", errors.New("syntax error "+pattern+" incomplete")
    }
    return gocode.String(), nil
}

func compileRenderer(lines []string, processed *list.List) (string, error) {
    gocode := new(bytes.Buffer)
    inside := false
    for i, l := range lines {
    }
}

