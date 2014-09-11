genresource
-----------
this is a utility is use to convert `../gokserver` to a resource.go file.
it allows, go binaries to contain resource extraction resource using this utility requires use of just the following snippet
```go
func UnpackResource() {
    orignal, err := hex.DecodeString(Resource)
    if err != nil {
        errExit(err, "")
    }
    unzip, err := gzip.NewReader(bytes.NewBuffer([]byte(orignal)))
    if err != nil {
        errExit(err, "")
    }
    untar := tar.NewReader(unzip)
    var fileContent bytes.Buffer
    for h, err := untar.Next(); err == nil; h, err = untar.Next() {
        fileContent.Reset()
        io.Copy(&fileContent, untar)
        var name string
        if h.Name[0] == '/' {
            name = h.Name[1:]
        } else {
            name = h.Name
        }
        if strings.Contains(name, "/") {
            lstIndx := strings.LastIndex(name, "/")
            if err := os.MkdirAll(name[0:lstIndx], 0744); err != nil {
                errExit(err, "")
            }
        }
        ioutil.WriteFile(name, fileContent.Bytes(), os.FileMode(h.Mode))
    }
}
```
