package main
import "fmt"
import "strings"
import "bytes"
import "errors"
import "container/list"
import "regexp"

var regsrch *regexp.Regexp
var regsrch2 *regexp.Regexp

func compile(code string) (string, string, error) {
    lines := strings.Split(code, "\n")
    processed := list.New()
    /*- create imports -*/
    imp, err := extract("<?goimp", lines, processed)
    if err != nil {
        return "", "", err
    }
    imports := ""
    if len(imp) != 0 && strings.TrimSpace(imp) != "" {
        imports = "import (\n"+imp+"\n)\n"
    }
    /*- create functions  -*/
    fn, err := extract("<?gofn", lines, processed)
    if err != nil {
        return "", "", err
    }
    funcs := processfnExtract(fn)
    /*- create use  -*/
    uses, err := extract("<?gouse", lines, processed)
    if err != nil {
        return "", "", err
    }
    /*- create renderer  -*/
    renderer := createRenderer(lines, processed)
    randName := "Render"+genRandName()
    fullcode := "package main\n"+imports+uses+funcs+"func GokRender"+randName+
                "(gok *Gok){\n"+renderer+"\n}"
    return fullcode, randName, err
}

func processfnExtract(s string) string {
    if (len(s) == 0) || (strings.TrimSpace(s) == "") {
        return ""
    }
    lns := strings.Split(s, "\n")
    lns[1] = "func "+lns[1]
    return strings.Join(lns, "\n")
}

type analyzed struct {
    ln int
    col1 int
    col2 int
}

func containsln(processed *list.List, ln int) (col1, col2 int) {
    for e := processed.Front(); e != nil; e = e.Next() {
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
            lnsAndCols.col2 = indx+2
            processed.PushBack(lnsAndCols)
        }
    }
    if inside {
        return "", errors.New("syntax error "+pattern+" incomplete")
    }
    return gocode.String(), nil
}

func createRenderer(lines []string, processed *list.List) string {
    gocode := new(bytes.Buffer)
    processlnStatic := false
    for i, l := range lines {
        lenl := len(l)
        col1, col2 := containsln(processed, i)
        var i1, i2 int
        if (col1 == 0) && (col2 == lenl) {
            continue
        } else if (col1 == -1) && (col2 == -1) {
            i1 = 0
            i2 = lenl
        } else if (col1 == 0) {
            i1 = col2
            i2 = lenl
        } else if (col2 == lenl) {
            i1 = 0
            i2 = col1
        }
        slice := l[i1:i2]
        gocode.Write(processln(slice, i+1, &processlnStatic))
    }
    return gocode.String()
}

func processln(s string, ln int, static *bool) []byte {
    if regsrch == nil {
        var err error
        regsrch, err = regexp.Compile("\\<\\?go\\s?[^(fn)(vars)(imp)]")
        if err != nil {
            return []byte("")
        }
    }
    gocode := new(bytes.Buffer)
    slen := len(s)
    for last := 0; last < slen; {
        slice := s[last:]
        indx1 := -1
        indx1End := 0
        temp := regsrch.FindIndex([]byte(slice))
        if temp != nil {
            indx1 = temp[0]
            indx1End = temp[1]
        }
        indx2 := strings.Index(slice, "?>")
        if (indx1 == -1) && (indx2 == -1) {
            /*- determine if this is html or go code -*/
            if (*static) {
                fmt.Fprintf(gocode, "//%d\n%s\n", ln, slice)
            } else {
                fmt.Fprintf(gocode, "//%d\ngok.Echo(\"%s\")\n", ln, buildStr(slice))
            }
            break
        } else if ((indx1 < indx2) && (indx1 != -1)) || (indx2 == -1) {
            if indx2 == -1 {
                indx2 = len(slice)
                *static = true
            }
            if strings.TrimSpace(slice[:indx1]) != "" {
                fmt.Fprintf(gocode, "//%d\ngok.Echo(\"%s\")\n", ln, buildStr(slice[:indx1]))
            }
            fmt.Fprintf(gocode, "//%d\n%s\n", ln, slice[indx1End:indx2])
            last += (indx2+2)
        } else {
            if indx1 == -1 {
                indx1 = len(slice)
                indx1End = 5
                (*static) = false
            }
            fmt.Fprintf(gocode, "//%d\n%s\n", ln, buildStr(slice[:indx2]))
            fmt.Fprintf(gocode, "//%d\ngok.Echo(\"%s\")\n", ln, slice[indx2+2:indx1])
            last += (indx1End)

        }
    }
    return gocode.Bytes()
}

func buildStr(s string) string {
    s = strings.Replace(s, "\n", "\\n", -1)
    s = strings.Replace(s, "\t", "\\t", -1)
    s = strings.Replace(s, "\r", "\\r", -1)
    s = strings.Replace(s, "\v", "\\v", -1)
    s = strings.Replace(s, "\f", "\\f", -1)
    s = strings.Replace(s, "\"", "\\\"", -1)
    return strings.TrimSpace(s)
}

