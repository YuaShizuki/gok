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
import "os"
import "encoding/hex"
import "math/rand"
import "time"


func pathExist(path string) (bool, os.FileInfo) {
    info, err :=  os.Stat(path);
    if err != nil {
        if os.IsNotExist(err) {
            return false, info;
        }
        panic("panic-> cannot determine if path exists");
    }
    return true, info;
}

func errExit(e error, msg string) {
    fmt.Println("error =>", e.Error());
    if msg != "" {
        fmt.Println(msg);
    }
    os.Exit(1);
}

func genRandName() string {
    rand.Seed(time.Now().UnixNano());
    num := rand.Uint32();
    str := []byte{ byte(num), byte(num >> 8), byte(num >> 16), byte(num >> 24) };
    return hex.EncodeToString(str);
}

