
package errpkg

import (
    "fmt"
    "runtime"
    "os"
)
var Errmsg string

func SetErrmsg(err string, values ...interface{}) {
    Errmsg = fmt.Sprintf(err, values...)
}

func InternError(err string) {
    _, fn, line, _ := runtime.Caller(1)
    fmt.Fprintf(os.Stderr, "%s:%d => %s\n", fn,
                line, err)
    os.Exit(1)
}