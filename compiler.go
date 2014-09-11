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
import "regexp"
import "strings"
import "bytes"
import "errors"

const(
    addfn = iota
    adduses = iota
    addrenderer = iota
    addimports = iota
    addajxfn = iota
    addunknown = iota
)
type processor func(string, int)(string, int)

var r1, r2, r3, r4, r5, r6, regend *regexp.Regexp
var ajaxFuncs map[string]string = nil
var regsearch map[*regexp.Regexp]processor = make(map[*regexp.Regexp]processor)

func compile(gokcode string) (string, string, map[string]string, error) {
    gokcodeLen := len(gokcode)
    //reset the map, populated during previous runs
    ajaxFuncs = nil
    ajaxFuncs = make(map[string]string)

    imports := new(bytes.Buffer)
    renderer := new(bytes.Buffer)
    uses := new(bytes.Buffer)
    funcs := new(bytes.Buffer)
    ajxfuncs := new(bytes.Buffer)

    if r1 == nil {
        r1,_ = regexp.Compile("\\<\\?gofn\\s")
        regsearch[r1] = processfn
        r2,_ = regexp.Compile("\\<\\?go\\s")
        regsearch[r2] = processgo
        r3,_ = regexp.Compile("\\<\\?gouse\\s")
        regsearch[r3] = processuse
        r4,_ = regexp.Compile("\\<\\?goimp\\s")
        regsearch[r4] = processimp
        r5,_ = regexp.Compile("\\<\\?go@fn\\s")
        regsearch[r5] = processajxfn
        regend,_ = regexp.Compile("\\?\\>")
    }

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
            lnoffset := getLineNum(gokcode[:(last+start[0])])
            if slice[start[1]-1] == '\n' {
                lnoffset++
            }
            echo := createEcho(slice[:start[0]], lnoffset)
            renderer.WriteString(echo)
            code, typ := v(slice[start[1]:end[0]],lnoffset,)
            switch typ {
                case addimports:
                    imports.WriteString(code)
                case addrenderer:
                    renderer.WriteString(code)
                case adduses:
                    uses.WriteString(code)
                case addfn:
                    funcs.WriteString(code)
                case addajxfn:
                    ajxfuncs.WriteString(code)
                case addunknown:
                    return "", "", nil, errors.New(code)
            }
            last += end[1]
            break
        }
    }
    rendererName := "Render"+genRandName()
    final := "package main\n"+imports.String()+uses.String()+funcs.String()+
            ajxfuncs.String()+"func "+rendererName+"(gok *Gok) {\n"+
            renderer.String()+"\n}\n"
    return final, rendererName, ajaxFuncs, nil
}

func getLineNum(code string) int {
    return strings.Count(code, "\n")+1
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

func processajxfn(code string, lnoff int) (string, int) {
    var out bytes.Buffer
    if r6 == nil {
        r6, _ = regexp.Compile("^\\s*[a-zA-Z0-9]+\\([a-zA-Z0-9]+\\s\\[\\]"+
        "string\\)\\s\\(\\[\\]string[,]\\serror\\)\\s*\\{(.|\\s)*\\}$")
    }
    if !r6.Match([]byte(code)) {
        err := fmt.Sprintf("go@fn inappropriate function on line %d",lnoff)
        return err, addunknown
    }
    indx := strings.Index(code, "(")
    if indx == -1 { return "fatal regexp failure", addunknown }
    fnName := strings.TrimSpace(code[:indx])
    newName := "ajx"+genRandName()
    //out := "func "+newName + code[indx:]
    out.WriteString(fmt.Sprintf("//%d\nfunc %s",lnoff, newName))
    lns := strings.Split(code[indx:], "\n")
    out.WriteString(lns[0]+"\n")
    if len(lns) > 1 {
        for i,l := range lns[1:] {
            out.WriteString(fmt.Sprintf("//%d\n", lnoff+i+1))
            out.WriteString(strings.TrimSpace(l)+"\n")
        }
    }
    out.WriteString("\n")
    ajaxFuncs[newName] = fnName
    return out.String(), addajxfn
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
    s = strings.Replace(s, "\\", "\\\\", -1)
    s = strings.Replace(s, "\n", "\\n", -1)
    s = strings.Replace(s, "\t", "\\t", -1)
    s = strings.Replace(s, "\r", "\\r", -1)
    s = strings.Replace(s, "\v", "\\v", -1)
    s = strings.Replace(s, "\f", "\\f", -1)
    s = strings.Replace(s, "\"", "\\\"", -1)
    return strings.TrimSpace(s)
}
