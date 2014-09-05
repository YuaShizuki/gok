package main

import "fmt"
import "strings"

type gok_ajxProcessor func([]string)([]string, error)

var gok_ajxRoutes map[string]gok_ajxProcessor = map[string]gok_ajxProcessor {
//<gok inject ajx routes>
}

func handleIfQuickAjx(gok *Gok) bool {
    if gok.ServerRequestMethod() != "POST" {
        return false
    }
    fn, ok := gok_ajxRoutes[gok.ServerSelf()]
    if !ok {
        return false
    }
    req := gok.Post("forgokqajxfn")
    result, err := fn(strings.Split(req,"[2577<--gokBoundry-->21501]"))
    if err != nil {
        fmt.Fprintf(gok.w, "");
    } else {
        fmt.Fprintf(gok.w, strings.Join(result, "[2577<--gokBoundry-->21501]"))
    }
    return true
}
