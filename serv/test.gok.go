package main
import "fmt"

func Render96ced45a(gok *Gok){

gok.Echo(" \n<html>\n    <body>\n        ");
for i := 0; i < 100; i++ {
gok.Echo("\n            <p>Hello World!</p>\n        ");
}
gok.Echo("\n    </body>\n</html>\n");
fmt.Println("responded to request => /test1")
gok.Echo("\n");
}
