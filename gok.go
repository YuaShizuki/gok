package main
import "fmt"
import "os"

func main() {
    if len(os.Args) != 2 {
        printUsage()
    }
    switch os.Args[1] {
        case "build":
            err := build(false);
            if err != nil {
                fmt.Println(err.Error());
            }
        case "run":
            fmt.Println("under construction")
        case "src":
            build(true)
        case "api":
            fmt.Println("under construction")
    }
}

