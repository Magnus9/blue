
package repl

import (
    "fmt"
    "os"
    "os/user"
    "bytes"
    "strings"
    "github.com/Magnus9/blue/parser"
    "github.com/Magnus9/blue/objects"
    "github.com/Magnus9/blue/blue"
    "github.com/chzyer/readline"
)

/*
 * The keywords that can start a blocked stmt. This
 * works really bad since newlines aint treated
 * very nicely in the parser (very strict). Something
 * better will come along later, the goal for now was
 * to get a repl up and running. Hack for now is to
 * use the newline replacement character '\'.

 * The goal strategy was to link the beginning keyword
 * of a statement onto a tokentype, and only recurse
 * the lines if the tokentype at the end of the stmt
 * was correct. Etc 'class' => NAME == Recurse.
 */
var kwds = map[string]struct{}{
    "def"  : struct{}{},
    "if"   : struct{}{},
    "while": struct{}{},
    "class": struct{}{},
}

func recurseLines(buf *bytes.Buffer) bool {
    for {
        line, err := readline.Line(">>> ")
        if err != nil {
            return false
        }
        if len(line) == 0 {
            continue
        }
        readline.AddHistory(line)

        buf.WriteByte('\n')
        buf.WriteString(line)
        /*
         * Maybe we have to recurse even further if we
         * find another stmt that yields a block.
         */
        list := strings.Fields(line)
        _, ok := kwds[list[0]]
        if ok {
            ok = recurseLines(buf)
            if !ok {
                return false
            }
            continue
        }
        /*
         * Check for the keyword that can end the block.
         * For now we only have 'end' for all blocks, so we
         * just have to check for that by iterating the list
         * we declared above.
         */
        for _, s := range list {
            if s == "end" {
                return true
            }
        }
    }
}

func readLine() (string, bool) {
    var buf bytes.Buffer
    for {
        line, err := readline.Line(">> ")
        if err != nil {
            goto errv
        }
        if len(line) == 0 {
            continue
        }
        readline.AddHistory(line)

        buf.WriteString(line)
        _, ok := kwds[strings.Fields(line)[0]]
        if ok {
            ok = recurseLines(&buf)
            if !ok {
                goto errv
            }
        }
        return buf.String(), true
    }
errv:
    return "", false
}

func Run(globals map[string]objects.BlObject) {
    defer func() {
        err := recover()
        if err != nil {
            fmt.Fprintf(os.Stderr, "%s\n", err)
            /*
             * Jump back into 'Run' again to process
             * the next stmt.
             */
            Run(globals)
        }
    }()
    for {
        program, ok := readLine()
        if !ok {
            break
        }
        ast := parser.ParseFromRepl("repl", program)
        runtime := blue.New("repl", ast, globals)
        runtime.Run()
    }
}

func Init() {
    us, err := user.Current()
    if err != nil {
        panic(err)
    }
    fmt.Println("Blue interactive shell")
    if len(us.HomeDir) == 0 {
        return
    }
    readline.SetHistoryPath(us.HomeDir + "/.blue_hist")
}