
package main
import "fmt"

func Render29d9e85a(gok *Gok){

gok.Echo(" \n<html>\n    <head>\n        <script>\n            function Write() { \n                body = document.getElementsByTagName(\"body\");\n                body[0].innerHTML += \"<p>Complete!</p>\";\n            }\n        </script>\n    </head>\n    <body onload=\"Write()\">\n        <p>\n            ");
for i := 0; i < 10; i++ {
gok.Echo(i,"<=>");
}
gok.Echo("\n        </p>\n    </body>\n</html>\n");
fmt.Println("responded to request =>"+gok.ServerSelf())
gok.Echo("\n");
}