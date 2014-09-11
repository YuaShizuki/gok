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
import "io/ioutil"
import "bytes"
import "archive/tar"
import "compress/gzip"
import "encoding/hex"
import "os"

func main() {
    if exist, _ := pathExist("resource.go"); exist {
        fmt.Println("resource.go already exists");
        return;
    }
    if len(os.Args) != 3 {
        fmt.Println("Usage: $genResource [dir] [package name]");
        return;
    }
    resourceFile := ""
    if dir,err := os.Getwd(); err != nil {
        errExit(err, "");
    } else {
        resourceFile = dir+"/resource.go";
    }
    if err := os.Chdir(os.Args[1]); err != nil {
        errExit(err, "");
    }
    buff := new(bytes.Buffer);
    tw := tar.NewWriter(buff);
    err := packDir("", tw);
    if err != nil {
        errExit(err, "");
    }
    if err := tw.Close(); err != nil {
        errExit(err, "");
    }
    finalBuff := new(bytes.Buffer);
    gw := gzip.NewWriter(finalBuff);
    gw.Write(buff.Bytes());
    gw.Close();
    //build a hex string for writing to a .go file
    hexStr := hex.EncodeToString(finalBuff.Bytes());
    f, err := os.Create(resourceFile);
    if err != nil {
        errExit(err, "");
    }
    f.Write([]byte("package "));
    f.Write([]byte(os.Args[2]+"\n"));
    f.Write([]byte("var Resource string = `"));
    f.Write([]byte(hexStr));
    f.Write([]byte("`;\n"));
    f.Close();
}

func packDir(dir string, w *tar.Writer) error {
    fileInfo, err := ioutil.ReadDir("."+dir);
    if err != nil {
        errExit(err, "");
    }
    for _, finfo := range fileInfo {
        if finfo.IsDir() {
            if err := packDir(dir+"/"+finfo.Name(), w); err != nil {
                errExit(err, "");
            }
            continue;
        }
        fileContent, err := ioutil.ReadFile("."+dir+"/"+finfo.Name());
        if err != nil {
            return err;
        }
        hdr := &tar.Header{
            Name:dir+"/"+finfo.Name(),
            Size:int64(len(fileContent)),
            Mode:int64(finfo.Mode()),
        };
        if err := w.WriteHeader(hdr); err != nil {
            return err;
        }
        if _, err := w.Write(fileContent); err != nil {
            return err;
        }
    }
    return nil;
}
