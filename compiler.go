package main
import "fmt"
import "strings"
import "io/ioutil"
import "bytes"
import "errors"

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
    ignoreindx := make(map[int]bool);
    header, err := extract("<?gohead", lines, ignoreindx)
    code, err := compileRenderer(lines, ignoreindx)
    return header, "", err
}

func extract(pattern string, lines []string) (string, error) {
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
            inside = true;
            if (lenl == plen) || (strings.TrimSpace(l[indx+plen:]) == "") {
                continue
            }
            fmt.Fprintf(gocode, "//%d\n%s\n",(i+1), strings.TrimSpace(l[indx+plen:]))
        } else {
            indx := strings.Index(l, "?>")
            if indx == -1 {
                fmt.Fprintf(gocode, "//%d\n%s\n", (i+1), strings.TrimSpace(l))
                continue
            }
            inside = false;
            if lenl == 2 {
                continue
            }
            if (indx == 0) || (strings.TrimSpace(l[:indx]) == "") {
                continue
            }
            fmt.Fprintf(gocode, "//%d\n%s\n", (i+1), strings.TrimSpace(l[:indx]))
        }
    }
    if inside {
        return "", errors.New("syntax error "+pattern+" incomplete")
    }
    return gocode.String(), nil
}

func compileRenderer(lines []string, skipNodes []string) (string, error) {
    gocode := new(bytes.Buffer)
    insideSkip := false
}

