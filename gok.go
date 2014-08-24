package main
import "fmt"
import "os"

func main() {
    if len(os.Args) != 2 {
        printUsage()
    }
    switch os.Args[1] {
        case "build":
            err := build();
            if err != nil {
                fmt.Println(err.Error());
            }
        case "run":
            fmt.Println("under construction")
        case "src":
            fmt.Println("under construction")
        case "api":
            fmt.Println("under construction")
    }
}

