package main
import "strings"
import "fmt"
import "bytes"
import "io/ioutil"

var gokJsCode string = `
function Actual(path, args) {
    boundry = "[2577<--gokBoundry-->21501]"
    pstData = ""
    for(var i=0; i<args.length; i++) {
        if ((typeof args[i] == "number") || (typeof args[i] == "string"))
            pstData += String(args[i]) + boundry 
    }
    if (pstData != "") {
        xhr = new XMLHttpRequest()
        xhr.open("POST", "/"+path, false)
        xhr.setRequestHeader("Content-type","application/x-www-form-urlencoded")
        xhr.send("forgokqajxfn="+pstData)
        return xhr.response.split("[2577<--gokBoundry-->21501]")
    }
}
function Gok() {
    this.init = true
}
`

var protJsCode string =
`Gok.prototype.%s = function(){ return Actual("%s", arguments) }
`

func buildGokJs(ajxfn map[string]string) {
    var out bytes.Buffer
    out.WriteString(gokJsCode)
    for k,v := range ajxfn {
        out.WriteString(fmt.Sprintf(protoJsCode, v, k))
    }
    out.WriteString("var gok = new Gok()")
    ioutil.WriteFile("gok.js", out.Bytes(), 0644)
}

func injectAjxRoutes(ajxroutes map[string]string) {
    var out bytes.Buffer
    file, err := ioutil.ReadFile("ajx.auto.go")
    if err != nil {
        errExit(err, "")
    }
    ajxroutes := strings.Split(string(file), "//<gok inject ajx routes>")
    if len(ajxroutes) != 2 {
        errExit(errors.New("corrupt ajx.auot.go file in resource"), "")
    }
    out.WriteString(ajxroutes[0])
    for k, v := range ajxroutes {
        out.WriteString("\""+k+"\":" + k + ",\n")
    }
    out.WriteString(ajxroutes[1])
    err = ioutil.WriteFile("ajx.auot.go", out.Bytes(), 0644)
    if err != nil {
        errExit(err, "")
    }
}
