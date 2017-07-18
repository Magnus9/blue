
package blue

import (
    "fmt"
    "bytes"
    "github.com/Magnus9/blue/errpkg"
    "github.com/Magnus9/blue/objects"
)

func genericTraceFunc(frame *objects.BlFrame) {
    var buf bytes.Buffer
    for f := frame; f != nil; f = f.Prev {
        str := fmt.Sprintf("in %s:%d\n   %s\n", f.Pathname,
                           f.Node.LineNum, f.Node.Line)
        buf.WriteString(str)
    }
    buf.WriteString(errpkg.Errmsg)

    panic(buf.String())
}