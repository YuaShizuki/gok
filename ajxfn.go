/*
- Copyright 2014 Keshav Bhide. All rights reserved.
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
import "strings"
import "bytes"
import "io/ioutil"
import "errors"

var gokJsCode string = `
function Actual(path, args) {
    boundry = "[2577<--gokBoundry-->21501]"
    pstData = ""
    for(var i=0; i<args.length; i++) {
        if ((typeof args[i] == "number") || (typeof args[i] == "string"))
            pstData += String(args[i]) + boundry 
    }
    xhr = new XMLHttpRequest()
    xhr.onload = function() { 
        args[args.length - 1](xhr.response.split("[2577<--gokBoundry-->21501]")) 
    }
    xhr.open("POST", "/"+path, true)
    xhr.setRequestHeader("Content-type","application/x-www-form-urlencoded")
    xhr.send("forgokqajxfn="+pstData)
}
function Gok() {
    this.init = true
}
`

var protoJsCode string =
`Gok.prototype.%s = function(){ Actual("%s", arguments) }
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

func injectAjxRoutes(ajxr map[string]string) {
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
    for k,_ := range ajxr {
        out.WriteString("\""+k+"\":" + k + ",\n")
    }
    out.WriteString(ajxroutes[1])
    err = ioutil.WriteFile("ajx.auto.go", out.Bytes(), 0644)
    if err != nil {
        errExit(err, "")
    }
}
