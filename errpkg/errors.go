
package errpkg

import (
    "fmt"
    "os"
    "bytes"
    "runtime"
    "path/filepath"
)
var Errmsg string

func SetErrmsg(err string, values ...interface{}) {
    Errmsg = fmt.Sprintf(err, values...)
}

func InternError(err string, values ...interface{}) {
    _, fn, line, _ := runtime.Caller(2)
    var buf bytes.Buffer
    buf.WriteString(fmt.Sprintf("%s:%d => %s\n",
                    filepath.Base(fn), line, err))
    fmt.Fprintf(os.Stderr, buf.String(), values...)
    // That exit status 1 is bleh.
    os.Exit(0)
}