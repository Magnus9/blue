
package main

import (
    "os"
    "fmt"
    "github.com/Magnus9/blue/parser"
    "github.com/Magnus9/blue/repl"
    "github.com/Magnus9/blue/objects"
    "github.com/Magnus9/blue/blue"
)

func main() {
    blue.Init()
    globals := make(map[string]objects.BlObject, 0)
    
    if len(os.Args) > 1 {
        pathname := os.Args[1]
        f, err := os.Open(pathname)
        if err != nil {
            fmt.Fprintf(os.Stderr, "%s\n", err)
            return
        }
        /*
        defer func() {
            err := recover()
            if err != nil {
                fmt.Println(err)
            }
        }()*/
        ast := parser.ParseFromFile(pathname, f)
        runtime := blue.New(pathname, ast, globals)
        runtime.Run()
    } else {
        repl.Init()
        repl.Run(globals)
    }
}