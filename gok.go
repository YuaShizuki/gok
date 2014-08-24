package main
import "fmt"

func main() {
    if len(os.Args) != 2 {
        printUsage();
    }
}

