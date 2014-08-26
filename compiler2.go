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
    final := new(bytes.Buffer)

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
    
}
