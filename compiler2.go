package main
import "fmt"
import "regexp"
import "io/ioutil"
import "strings"


const(
    addfn = iota
    adduses = iota
    addrenderer = iota
    addimports = iota
)

func main() {
    f, err := ioutil.ReadFile("index.gok")
    if err != nil {
        fmt.Println("error reading file")
        return
    }
    fmt.Println(compile(string(f)))
}

func compile(gokcode string) string {
    gokcodeLen := len(gokcode)

    imports := new(bytes.Buffer)
    renderer := new(bytes.Buffer)
    uses := new(bytes.Buffer)
    funcs := new(byes.Buffer)

    regsearch map[string]func(string, int) (string, int)
    regsearch[regexp.Compile("\\<\\?gofn\\s")] = processfn
    regsearch[regexp.Compile("\\<\\?go\\s")] = processgo
    regsearch[regexp.Compile("\\<\\?gouse\\s")] = processuse
    regsearch[regexp.Compile("\\<\\?goimp\\s")] = processimp

    regend := regexp.Compile("\\?\\>")

    for last := 0; last < gokcodeLen {
        slice := gokcode[last:]
        end := regend.FindIndex(slice)
        if end == nil {
            echo := createEcho(slice[:], getLineNum(gokcode[:last]))
            renderer.Write(echo)
        }
        for k, v := range regsearch {
            start := k.FindIndex(slice)
            if (start == nil) || (start[1] > end[0]) {
                continue
            }
            lnoffset := getLineNum(gokcode[:last])
            echo := createEcho(slice[:start[0]], lnoffset)
            renderer.Write(echo)
            code, typ := []byte(v(slice[start[1]:end[0]],lnoffset))
            switch typ {
                case addImports:
                    imports.Write(code)
                case addrenderer:
                    renderer.Write(code)
                case adduses:
                    uses.Write(code)
                case addfn:
                    funcs.Write(code)
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
    out.WriteString(lns[0]+"\n")
    for i := 1; i < llns; i++ {
        out.WriteString(fmt.Sprintf("//%d\n%s\n", lnoff+i, lns[i]))
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
        out.WriteString(fmt.Sprintf("//%d\n%s\n",i+lnoff,s))
    }
    return out.String(), addrenderer
}

func processuse(code string, lnoff int) (string, int) {
    out := new(bytes.Buffer)
    lns := string.Split(code, "\n")
    if len(lns) == 0 {
        return "", adduses
    }
    for i,s := range lns {
        out.WriteString(fmt.Sprintf("//%d\n%s\n", lnoff+i, s))
    }
    return out.String(), adduses
}

func procesimp(code string, lnoff int) {
    lns := strings.Split(code)
    if len(lns) == 0 {
        return "", addimports
    }
    out := new(bytes.Buffer)
    out.WriteString("import(\n")
    for i, s := range lns {
        out.WriteString(fmt.Sprintf("//%d\n%s\n", lnoff+i, s))
    }
    out.WriteString(")\n")
    return out.String(), addimports
}

