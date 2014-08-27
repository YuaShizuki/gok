package main
import "fmt"
import "regexp"
import "io/ioutil"
import "strings"
import "bytes"

const(
    addfn = iota
    adduses = iota
    addrenderer = iota
    addimports = iota
)
type processor func(string, int)(string, int)

func compile(gokcode string) string {
    gokcodeLen := len(gokcode)

    imports := new(bytes.Buffer)
    renderer := new(bytes.Buffer)
    uses := new(bytes.Buffer)
    funcs := new(bytes.Buffer)

    regsearch := make(map[*regexp.Regexp]processor)
    r1,_ := regexp.Compile("\\<\\?gofn\\s")
    regsearch[r1] = processfn
    r2,_ := regexp.Compile("\\<\\?go\\s")
    regsearch[r2] = processgo
    r3,_ := regexp.Compile("\\<\\?gouse\\s")
    regsearch[r3] = processuse
    r4,_:= regexp.Compile("\\<\\?goimp\\s")
    regsearch[r4] = processimp

    regend,_ := regexp.Compile("\\?\\>")

    for last := 0; last < gokcodeLen; {
        slice := gokcode[last:]
        end := regend.FindIndex([]byte(slice))
        if end == nil {
            echo := createEcho(slice[:], getLineNum(gokcode[:last]))
            renderer.WriteString(echo)
            break
        }
        for k, v := range regsearch {
            start := k.FindIndex([]byte(slice))
            if (start == nil) || (start[1] > end[0]) {
                continue
            }
            lnoffset := getLineNum(gokcode[:last])
            echo := createEcho(slice[:start[0]], lnoffset)
            renderer.WriteString(echo)
            code, typ := v(slice[start[1]:end[0]],lnoffset)
            switch typ {
                case addimports:
                    imports.WriteString(code)
                case addrenderer:
                    renderer.WriteString(code)
                case adduses:
                    uses.WriteString(code)
                case addfn:
                    funcs.WriteString(code)
            }
            last += end[1]
            break
        }
    }
    return "package main\n"+imports.String()+uses.String()+funcs.String()+
        "func Render(gok *Gok){\n"+renderer.String()+"\n}\n"
}

func getLineNum(code string) int {
    return strings.Count(code, "\n")
}


func processfn(code string, lnoff int) (string, int) {
    out := new(bytes.Buffer)
    out.WriteString(fmt.Sprintf("//%d\nfunc ", lnoff))
    lns := strings.Split(code, "\n")
    llns := len(lns)
    if llns == 0 {
        return "", addfn
    }
    out.WriteString(strings.TrimSpace(lns[0])+"\n")
    for i := 1; i < llns; i++ {
        s := strings.TrimSpace(lns[i])
        out.WriteString(fmt.Sprintf("//%d\n%s\n", lnoff+i, s))
    }
    return out.String(), addfn
}

func processgo(code string, lnoff int) (string, int) {
    out := new(bytes.Buffer)
    lns := strings.Split(code, "\n")
    if len(lns) == 0 {
        return "", addrenderer
    }
    for i,s := range lns {
        s := strings.TrimSpace(s)
        out.WriteString(fmt.Sprintf("//%d\n%s\n",i+lnoff, s))
    }
    return out.String(), addrenderer
}

func processuse(code string, lnoff int) (string, int) {
    out := new(bytes.Buffer)
    lns := strings.Split(code, "\n")
    if len(lns) == 0 {
        return "", adduses
    }
    for i,s := range lns {
        s := strings.TrimSpace(s)
        out.WriteString(fmt.Sprintf("//%d\n%s\n", lnoff+i, s))
    }
    return out.String(), adduses
}

func processimp(code string, lnoff int) (string, int) {
    lns := strings.Split(code, "\n")
    if len(lns) == 0 {
        return "", addimports
    }
    out := new(bytes.Buffer)
    out.WriteString("import(\n")
    for i, s := range lns {
        s = strings.TrimSpace(s)
        out.WriteString(fmt.Sprintf("//%d\n%s\n", lnoff+i, s))
    }
    out.WriteString(")\n")
    return out.String(), addimports
}

func createEcho(code string, lnoff int) string {
    formatedStr := formateStr(code)
    if formatedStr != "" {
        return fmt.Sprintf("//%d\ngok.Echo(\"%s\")\n", lnoff, formatedStr)
    }
    return ""
}

func formateStr(s string) string {
    if strings.TrimSpace(s) == "" {
        return ""
    }
    s = strings.Replace(s, "\n", "\\n", -1)
    s = strings.Replace(s, "\t", "\\t", -1)
    s = strings.Replace(s, "\r", "\\r", -1)
    s = strings.Replace(s, "\v", "\\v", -1)
    s = strings.Replace(s, "\f", "\\f", -1)
    s = strings.Replace(s, "\"", "\\\"", -1)
    return strings.TrimSpace(s)
}
